package dev.hgmx.starter

import androidx.test.ext.junit.runners.AndroidJUnit4
import androidx.test.platform.app.InstrumentationRegistry
import com.connectrpc.Code
import com.connectrpc.ConnectErrorDetail
import com.connectrpc.ConnectException
import com.connectrpc.ResponseMessage
import com.connectrpc.extensions.GoogleJavaLiteProtobufStrategy
import com.google.protobuf.MessageLite
import hgmx.app.v1.UpdateRequired
import hgmx.auth.v1.Account
import hgmx.auth.v1.ClientKind
import hgmx.auth.v1.GetAccountResponse
import hgmx.auth.v1.SessionGrant
import hgmx.auth.v1.SignInAgain
import hgmx.auth.v1.VerifyCodeResponse
import java.io.File
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.CompletableDeferred
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.async
import kotlinx.coroutines.awaitCancellation
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking
import okio.ByteString.Companion.toByteString
import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertNull
import org.junit.Assert.assertThrows
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test
import org.junit.runner.RunWith
import hgmx.auth.v1.Session as Device

@RunWith(AndroidJUnit4::class)
class SessionTest {
    private val context = InstrumentationRegistry.getInstrumentation().targetContext
    private val dir = File(context.cacheDir, "session-test")
    private val tokens = TokenStore(dir)
    private val scope = CoroutineScope(SupervisorJob() + Dispatchers.Default)
    private val session = Session(Api("https://example.com", "test-android", 1), tokens, scope)
    private val token = "a".repeat(43)

    @Before
    fun reset() {
        dir.deleteRecursively()
        dir.mkdirs()
    }

    @Test
    fun tokenIsEncryptedAtRest() {
        tokens.write(token)
        assertEquals(token, tokens.read())
        assertFalse(dir.walk().filter { it.isFile }.any { it.readText().contains(token) })
        tokens.clear()
        assertNull(tokens.read())
    }

    @Test
    fun lateGrantAfterSignOutIsDropped() = runBlocking {
        val sent = CompletableDeferred<Unit>()
        val reply = CompletableDeferred<ResponseMessage<VerifyCodeResponse>>()
        val signIn = async {
            session.signIn({ sent.complete(Unit); reply.await() }) { it.grant }
        }
        sent.await()
        session.signOut()
        session.gate.first { it is Gate.SignedOut }
        reply.complete(success(VerifyCodeResponse.newBuilder().setGrant(grant()).build()))
        assertThrows(CancellationException::class.java) { runBlocking { signIn.await() } }
        assertTrue(session.gate.value is Gate.SignedOut)
        assertNull(tokens.read())
    }

    @Test
    fun grantIsStoredAndReplacedOnReauth() = runBlocking {
        session.signIn({ success(VerifyCodeResponse.newBuilder().setGrant(grant()).build()) }) { it.grant }
        assertEquals(token, tokens.read())
        assertTrue(session.gate.value is Gate.SignedIn)
    }

    @Test
    fun unconfirmedSignOutSaysSo() = runBlocking {
        val offline = Session(Api("http://localhost:18559", "test-android", 1), tokens, scope)
        offline.signIn({ success(VerifyCodeResponse.newBuilder().setGrant(grant()).build()) }) { it.grant }
        offline.signOut()
        val notice = offline.gate.first { it is Gate.SignedOut && it.notice != null }
            .let { (it as Gate.SignedOut).notice.orEmpty() }
        assertNull(tokens.read())
        assertTrue(notice, notice.contains("didn't confirm"))
        assertTrue(notice, notice.contains("http://localhost:18559/auth/account"))
    }

    @Test
    fun offlineRecheckKeepsTheGate() = runBlocking {
        val offline = Session(Api("http://localhost:18559", "test-android", 1), tokens, scope)
        offline.signIn({ success(VerifyCodeResponse.newBuilder().setGrant(grant()).build()) }) { it.grant }
        val gate = offline.gate.value
        offline.recheck()
        assertEquals(gate, offline.gate.value)
    }

    @Test
    fun unauthenticatedEndsTheSession() = runBlocking {
        session.signIn({ success(VerifyCodeResponse.newBuilder().setGrant(grant()).build()) }) { it.grant }
        val failure = ConnectException(Code.UNAUTHENTICATED, "Your session ended. Sign in again.")
        val problem = assertThrows(Problem::class.java) {
            runBlocking { session.call<GetAccountResponse> { ResponseMessage.Failure(failure, emptyMap(), emptyMap()) } }
        }
        assertEquals("Your session ended. Sign in again.", problem.message)
        assertEquals(Gate.SignedOut("Your session ended. Sign in again."), session.gate.value)
        assertNull(tokens.read())
    }

    @Test
    fun typedDetailsSteerTheApp() = runBlocking {
        val reauth = assertThrows(Problem::class.java) {
            runBlocking { session.call<GetAccountResponse> { failure("hgmx.auth.v1.SignInAgain", SignInAgain.newBuilder().setMessage("Sign in again to continue.").build()) } }
        }
        assertTrue(reauth.reauth)
        assertEquals(Gate.Starting, session.gate.value)

        val update = UpdateRequired.newBuilder().setMinimumBuild(7).setStoreUrl("https://play.google.com/store").setMessage("Update to keep going.").build()
        assertThrows(Problem::class.java) {
            runBlocking { session.call<GetAccountResponse> { failure("type.googleapis.com/hgmx.app.v1.UpdateRequired", update) } }
        }
        assertEquals(Gate.Update("https://play.google.com/store", "Update to keep going."), session.gate.value)
    }

    @Test
    fun cancelledCallLeavesStateAlone() = runBlocking {
        val call = async { session.call<GetAccountResponse> { awaitCancellation() } }
        call.cancel()
        assertTrue(runCatching { call.await() }.exceptionOrNull() is CancellationException)
        assertEquals(Gate.Starting, session.gate.value)
    }

    @Test
    fun unknownEnumsParse() {
        val bytes = byteArrayOf(0x18, 0x63)
        assertEquals(ClientKind.UNRECOGNIZED, Device.parseFrom(bytes).kind)
        assertEquals(99, Device.parseFrom(bytes).kindValue)
    }

    @Test
    fun originMustBeExact() {
        for (bad in listOf("http://example.com", "https://example.com/api", "https://user@example.com", "https://example.com?x=1")) {
            assertThrows(IllegalArgumentException::class.java) { Api(bad, "test-android", 1) }
        }
        Api("http://localhost:18550", "test-android", 1)
    }

    private fun grant() = SessionGrant.newBuilder()
        .setToken(token)
        .setAccount(Account.newBuilder().setId("a1").setUsername("ada"))
        .build()

    private fun <T> success(message: T) = ResponseMessage.Success(message, emptyMap(), emptyMap())

    private fun <T> failure(type: String, detail: MessageLite): ResponseMessage<T> {
        val e = ConnectException(Code.FAILED_PRECONDITION, "precondition")
            .withErrorDetails(GoogleJavaLiteProtobufStrategy().errorDetailParser(), listOf(ConnectErrorDetail(type, detail.toByteArray().toByteString())))
        return ResponseMessage.Failure(e, emptyMap(), emptyMap())
    }
}
