package dev.hgmx.starter

import android.os.Build
import com.connectrpc.Code
import com.connectrpc.ConnectException
import com.connectrpc.Headers
import com.connectrpc.ResponseMessage
import hgmx.app.v1.GetConfigurationRequest
import hgmx.app.v1.GetConfigurationResponse
import hgmx.app.v1.UpdateRequired
import hgmx.auth.v1.Account
import hgmx.auth.v1.DeleteAccountRequest
import hgmx.auth.v1.GetAccountRequest
import hgmx.auth.v1.SessionGrant
import hgmx.auth.v1.SignInAgain
import hgmx.auth.v1.SignOutRequest
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import kotlinx.coroutines.sync.Mutex
import kotlinx.coroutines.sync.withLock
import kotlinx.coroutines.withContext

sealed interface Gate {
    data object Starting : Gate
    data class Unreachable(val problem: Problem) : Gate
    data class Update(val storeUrl: String, val message: String) : Gate
    data class SignedOut(val notice: String? = null) : Gate
    data class SignedIn(val account: Account) : Gate
}

class Problem(message: String, val reauth: Boolean = false) : Exception(message)

class Session(
    val api: Api,
    private val tokens: TokenStore,
    private val scope: CoroutineScope,
) {
    private val _gate = MutableStateFlow<Gate>(Gate.Starting)
    val gate: StateFlow<Gate> = _gate.asStateFlow()
    private val _warning = MutableStateFlow<String?>(null)
    val warning: StateFlow<String?> = _warning.asStateFlow()

    private val lock = Mutex()
    private val checking = Mutex()
    private var token: String? = null
    private var generation = 0L
    private var methods = emptySet<String>()

    val label: String = Build.MODEL

    fun offers(method: String) = methods.any { it.endsWith("/$method") }

    suspend fun start() {
        _gate.value = Gate.Starting
        try {
            val config = configure()
            if (config.updateRequired) return
            methods = config.methodsList.toSet()
            val stored = withContext(Dispatchers.IO) { tokens.read() }
            lock.withLock {
                token = stored
                advance()
            }
            if (stored == null) {
                _gate.value = Gate.SignedOut()
                return
            }
            exchange({ api.auth.getAccount(GetAccountRequest.getDefaultInstance(), it) }) {
                _gate.value = Gate.SignedIn(it.account)
            }
        } catch (p: Problem) {
            if (_gate.value == Gate.Starting) _gate.value = Gate.Unreachable(p)
        }
    }

    suspend fun recheck() {
        val gate = _gate.value
        if (gate == Gate.Starting || gate is Gate.Update || !checking.tryLock()) return
        try {
            configure()
        } catch (_: Problem) {
        } finally {
            checking.unlock()
        }
    }

    suspend fun <T> call(rpc: suspend (Headers) -> ResponseMessage<T>): T = exchange(rpc)

    suspend fun <T> signIn(rpc: suspend (Headers) -> ResponseMessage<T>, grant: (T) -> SessionGrant) {
        var previous: String? = null
        exchange(rpc) {
            val g = grant(it)
            withContext(Dispatchers.IO) { tokens.write(g.token) }
            previous = token
            token = g.token
            advance()
            _gate.value = Gate.SignedIn(g.account)
        }
        previous?.let { scope.launch { revoke(it) } }
    }

    suspend fun <T> update(rpc: suspend (Headers) -> ResponseMessage<T>, account: (T) -> Account): Account =
        account(exchange(rpc) { _gate.value = Gate.SignedIn(account(it)) })

    suspend fun deleteAccount() {
        exchange({ api.auth.deleteAccount(DeleteAccountRequest.newBuilder().setConfirm(true).build(), it) }) {
            end(null)
        }
    }

    fun signOut() {
        scope.launch {
            val (old, sent) = lock.withLock {
                val old = token
                end(null)
                old to generation
            }
            if (old == null || revoke(old)) return@launch
            lock.withLock {
                if (sent == generation) {
                    _warning.value = "You're signed out on this device, but the server didn't confirm it. " +
                        "To end that session, sign in at ${api.origin}/auth/account and sign out $label under Devices."
                }
            }
        }
    }

    private suspend fun revoke(old: String): Boolean =
        when (val response = api.auth.signOut(SignOutRequest.getDefaultInstance(), api.headers(old))) {
            is ResponseMessage.Success -> true
            is ResponseMessage.Failure -> response.cause.code == Code.UNAUTHENTICATED
        }

    private suspend fun end(notice: String?) {
        withContext(Dispatchers.IO) { tokens.clear() }
        token = null
        advance()
        _gate.value = Gate.SignedOut(notice)
    }

    private fun advance() {
        generation++
        _warning.value = null
    }

    private suspend fun configure(): GetConfigurationResponse =
        exchange({ api.app.getConfiguration(GetConfigurationRequest.getDefaultInstance(), it) }) {
            if (it.updateRequired) _gate.value = Gate.Update(it.storeUrl, "Update the app to keep going.")
        }

    private suspend fun <T> exchange(
        rpc: suspend (Headers) -> ResponseMessage<T>,
        commit: suspend (T) -> Unit = {},
    ): T {
        val (sent, headers) = lock.withLock { generation to api.headers(token) }
        val response = rpc(headers)
        return lock.withLock {
            if (sent != generation) throw CancellationException("The session changed during this call.")
            when (response) {
                is ResponseMessage.Success -> response.message.also { commit(it) }
                is ResponseMessage.Failure -> throw problem(response.cause, signedIn = headers.containsKey("Authorization"))
            }
        }
    }

    private suspend fun problem(e: ConnectException, signedIn: Boolean): Problem {
        e.detail("hgmx.app.v1.UpdateRequired", UpdateRequired::parseFrom)?.let {
            _gate.value = Gate.Update(it.storeUrl, it.message)
            return Problem(it.message)
        }
        e.detail("hgmx.auth.v1.SignInAgain", SignInAgain::parseFrom)?.let {
            return Problem(it.message, reauth = true)
        }
        val message = when {
            e.exception != null -> "Couldn't reach ${api.origin}. Check your connection."
            else -> e.message ?: "Something went wrong. Try again."
        }
        if (e.code == Code.UNAUTHENTICATED && signedIn) end(message)
        return Problem(message)
    }
}

private fun <T> ConnectException.detail(type: String, parse: (ByteArray) -> T): T? =
    details.firstOrNull { it.type.substringAfterLast('/') == type }?.let { parse(it.payload.toByteArray()) }
