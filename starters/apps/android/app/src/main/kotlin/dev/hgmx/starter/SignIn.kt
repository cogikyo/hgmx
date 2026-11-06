package dev.hgmx.starter

import android.content.Context
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.Button
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp
import androidx.lifecycle.ViewModel
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.viewModelScope
import androidx.lifecycle.viewmodel.compose.viewModel
import com.google.protobuf.ByteString
import hgmx.auth.v1.BeginGoogleSigninRequest
import hgmx.auth.v1.BeginPasskeySigninRequest
import hgmx.auth.v1.FinishGoogleSigninRequest
import hgmx.auth.v1.FinishPasskeySigninRequest
import hgmx.auth.v1.SendCodeRequest
import hgmx.auth.v1.VerifyCodeRequest
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch

data class SignInState(
    val email: String = "",
    val code: String = "",
    val sent: Boolean = false,
    val busy: Boolean = false,
    val problem: Problem? = null,
    val done: Boolean = false,
)

class SignInViewModel(private val session: Session) : ViewModel() {
    private val api = session.api
    private val _state = MutableStateFlow(SignInState())
    val state = _state.asStateFlow()

    val codes = session.offers("SendCode")
    val passkeys = session.offers("BeginPasskeySignin")
    val google = session.offers("BeginGoogleSignin") && BuildConfig.GOOGLE_CLIENT_ID.isNotEmpty()

    fun email(value: String) = _state.update { it.copy(email = value, sent = false, code = "") }

    fun code(value: String) = _state.update { it.copy(code = value) }

    fun sendCode() = act {
        val email = _state.value.email.trim()
        session.call { api.auth.sendCode(SendCodeRequest.newBuilder().setEmail(email).build(), it) }
        _state.update { it.copy(sent = true) }
    }

    fun verify() = act {
        val s = _state.value
        val request = VerifyCodeRequest.newBuilder()
            .setEmail(s.email.trim())
            .setCode(s.code.trim())
            .setLabel(session.label)
            .build()
        session.signIn({ api.auth.verifyCode(request, it) }) { it.grant }
        _state.update { it.copy(done = true) }
    }

    fun passkey(context: Context) = act {
        val begin = session.call { api.auth.beginPasskeySignin(BeginPasskeySigninRequest.getDefaultInstance(), it) }
        val response = Credentials.passkey(context, begin.options.publicKeyOptionsJson.toStringUtf8())
        val request = FinishPasskeySigninRequest.newBuilder()
            .setCeremony(begin.options.ceremony)
            .setCredentialJson(ByteString.copyFromUtf8(response))
            .setLabel(session.label)
            .build()
        session.signIn({ api.auth.finishPasskeySignin(request, it) }) { it.grant }
        _state.update { it.copy(done = true) }
    }

    fun google(context: Context) = act {
        val flow = session.call { api.auth.beginGoogleSignin(BeginGoogleSigninRequest.getDefaultInstance(), it) }.flow
        val token = Credentials.google(context, BuildConfig.GOOGLE_CLIENT_ID, flow.nonce)
        val request = FinishGoogleSigninRequest.newBuilder()
            .setFlow(flow.flow)
            .setIdToken(token)
            .setLabel(session.label)
            .build()
        session.signIn({ api.auth.finishGoogleSignin(request, it) }) { it.grant }
        _state.update { it.copy(done = true) }
    }

    private fun act(action: suspend () -> Unit) {
        if (_state.value.busy) return
        _state.update { it.copy(busy = true, problem = null) }
        viewModelScope.launch {
            try {
                action()
            } catch (p: Problem) {
                _state.update { it.copy(problem = p) }
            } finally {
                _state.update { it.copy(busy = false) }
            }
        }
    }
}

@Composable
fun SignInScreen(session: Session, notice: String?, onDone: () -> Unit) {
    val vm = viewModel { SignInViewModel(session) }
    val state by vm.state.collectAsStateWithLifecycle()
    val context = LocalContext.current
    LaunchedEffect(state.done) {
        if (state.done) onDone()
    }
    Column(
        Modifier.fillMaxSize().padding(24.dp),
        verticalArrangement = Arrangement.spacedBy(12.dp, Alignment.CenterVertically),
    ) {
        Text("Sign in", style = MaterialTheme.typography.headlineMedium)
        notice?.let { Text(it) }
        if (vm.passkeys) {
            Button(onClick = { vm.passkey(context) }, enabled = !state.busy, modifier = Modifier.fillMaxWidth()) {
                Text("Sign in with a passkey")
            }
        }
        if (vm.google) {
            OutlinedButton(onClick = { vm.google(context) }, enabled = !state.busy, modifier = Modifier.fillMaxWidth()) {
                Text("Sign in with Google")
            }
        }
        if (vm.codes) {
            HorizontalDivider()
            OutlinedTextField(
                value = state.email,
                onValueChange = vm::email,
                label = { Text("Email") },
                singleLine = true,
                keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Email),
                modifier = Modifier.fillMaxWidth(),
            )
            if (state.sent) {
                OutlinedTextField(
                    value = state.code,
                    onValueChange = vm::code,
                    label = { Text("Code from your email") },
                    singleLine = true,
                    keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.NumberPassword),
                    modifier = Modifier.fillMaxWidth(),
                )
                Button(onClick = vm::verify, enabled = !state.busy && state.code.isNotBlank(), modifier = Modifier.fillMaxWidth()) {
                    Text("Sign in")
                }
            } else {
                OutlinedButton(onClick = vm::sendCode, enabled = !state.busy && state.email.isNotBlank(), modifier = Modifier.fillMaxWidth()) {
                    Text("Email me a code")
                }
            }
        }
        state.problem?.let { Text(it.message.orEmpty(), color = MaterialTheme.colorScheme.error) }
    }
}
