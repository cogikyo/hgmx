package dev.hgmx.starter

import android.content.Context
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Button
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.ListItem
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp
import androidx.lifecycle.ViewModel
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.viewModelScope
import androidx.lifecycle.viewmodel.compose.viewModel
import com.google.protobuf.ByteString
import com.google.protobuf.Timestamp
import hgmx.auth.v1.BeginGoogleLinkRequest
import hgmx.auth.v1.BeginPasskeyAddRequest
import hgmx.auth.v1.ChangeEmailRequest
import hgmx.auth.v1.FinishGoogleLinkRequest
import hgmx.auth.v1.FinishPasskeyAddRequest
import hgmx.auth.v1.Passkey
import hgmx.auth.v1.Provider
import hgmx.auth.v1.RemovePasskeyRequest
import hgmx.auth.v1.RenamePasskeyRequest
import hgmx.auth.v1.SetUsernameRequest
import hgmx.auth.v1.UnlinkGoogleRequest
import java.time.Instant
import java.time.ZoneId
import java.time.format.DateTimeFormatter
import java.time.format.FormatStyle
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch

data class AccountState(
    val busy: Boolean = false,
    val problem: Problem? = null,
    val notice: String? = null,
)

class AccountViewModel(private val session: Session) : ViewModel() {
    private val api = session.api
    private val _state = MutableStateFlow(AccountState())
    val state = _state.asStateFlow()

    fun username(name: String) = act {
        session.update({ api.auth.setUsername(SetUsernameRequest.newBuilder().setUsername(name).build(), it) }) { it.account }
    }

    fun email(email: String) = act {
        session.call { api.auth.changeEmail(ChangeEmailRequest.newBuilder().setEmail(email).build(), it) }
        _state.update { it.copy(notice = "Open the link we sent to $email to finish the change.") }
    }

    fun addPasskey(context: Context) = act {
        val begin = session.call { api.auth.beginPasskeyAdd(BeginPasskeyAddRequest.getDefaultInstance(), it) }
        val response = Credentials.createPasskey(context, begin.options.publicKeyOptionsJson.toStringUtf8())
        val request = FinishPasskeyAddRequest.newBuilder()
            .setCeremony(begin.options.ceremony)
            .setCredentialJson(ByteString.copyFromUtf8(response))
            .build()
        session.update({ api.auth.finishPasskeyAdd(request, it) }) { it.account }
    }

    fun rename(passkey: Passkey, name: String) = act {
        val request = RenamePasskeyRequest.newBuilder().setPasskeyId(passkey.id).setName(name).build()
        session.update({ api.auth.renamePasskey(request, it) }) { it.account }
    }

    fun remove(passkey: Passkey) = act {
        val request = RemovePasskeyRequest.newBuilder().setPasskeyId(passkey.id).build()
        session.update({ api.auth.removePasskey(request, it) }) { it.account }
    }

    fun linkGoogle(context: Context) = act {
        val flow = session.call { api.auth.beginGoogleLink(BeginGoogleLinkRequest.getDefaultInstance(), it) }.flow
        val token = Credentials.google(context, BuildConfig.GOOGLE_CLIENT_ID, flow.nonce)
        val request = FinishGoogleLinkRequest.newBuilder().setFlow(flow.flow).setIdToken(token).build()
        session.update({ api.auth.finishGoogleLink(request, it) }) { it.account }
    }

    fun unlinkGoogle() = act {
        session.update({ api.auth.unlinkGoogle(UnlinkGoogleRequest.getDefaultInstance(), it) }) { it.account }
    }

    fun delete() = act { session.deleteAccount() }

    private fun act(action: suspend () -> Unit) {
        if (_state.value.busy) return
        _state.update { it.copy(busy = true, problem = null, notice = null) }
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

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun AccountScreen(session: Session, open: (Route) -> Unit, back: () -> Unit) {
    val vm = viewModel { AccountViewModel(session) }
    val state by vm.state.collectAsStateWithLifecycle()
    val gate by session.gate.collectAsStateWithLifecycle()
    val account = (gate as? Gate.SignedIn)?.account ?: return
    val context = LocalContext.current
    val caps = account.capabilities
    var username by rememberSaveable { mutableStateOf(account.username) }
    var email by rememberSaveable { mutableStateOf("") }
    var renaming by rememberSaveable { mutableStateOf<String?>(null) }
    var deleting by rememberSaveable { mutableStateOf(false) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Account") },
                navigationIcon = { TextButton(onClick = back) { Text("Back") } },
            )
        },
    ) { padding ->
        LazyColumn(
            Modifier.padding(padding).padding(horizontal = 16.dp),
            verticalArrangement = Arrangement.spacedBy(12.dp),
        ) {
            item {
                ProblemText(state.problem, open)
                state.notice?.let { Text(it) }
            }
            item {
                Text(account.displayName, style = MaterialTheme.typography.titleLarge)
                if (account.hasEmail()) Text(account.email)
            }
            item {
                OutlinedTextField(username, { username = it }, label = { Text("Username") }, singleLine = true, modifier = Modifier.fillMaxWidth())
                Button(onClick = { vm.username(username.trim()) }, enabled = !state.busy && username.trim() != account.username) {
                    Text("Save username")
                }
            }
            if (caps.changeEmail) {
                item {
                    OutlinedTextField(
                        email, { email = it }, label = { Text("New email") }, singleLine = true,
                        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Email),
                        modifier = Modifier.fillMaxWidth(),
                    )
                    Button(onClick = { vm.email(email.trim()) }, enabled = !state.busy && email.isNotBlank()) { Text("Change email") }
                }
            }
            item {
                HorizontalDivider()
                Text("Passkeys", style = MaterialTheme.typography.titleMedium)
            }
            items(account.passkeysList, key = { it.id }) { passkey ->
                ListItem(
                    headlineContent = { Text(passkey.name) },
                    supportingContent = { Text(passkeyDetail(passkey)) },
                    trailingContent = {
                        Row {
                            TextButton(onClick = { renaming = passkey.id }, enabled = !state.busy) { Text("Rename") }
                            if (passkey.removable) {
                                TextButton(onClick = { vm.remove(passkey) }, enabled = !state.busy) { Text("Remove") }
                            }
                        }
                    },
                )
            }
            item {
                OutlinedButton(onClick = { vm.addPasskey(context) }, enabled = !state.busy) { Text("Add a passkey") }
            }
            if (caps.google) {
                item {
                    HorizontalDivider()
                    val google = account.identitiesList.firstOrNull { it.provider == Provider.PROVIDER_GOOGLE }
                    when {
                        google == null -> OutlinedButton(onClick = { vm.linkGoogle(context) }, enabled = !state.busy && BuildConfig.GOOGLE_CLIENT_ID.isNotEmpty()) {
                            Text("Link Google")
                        }
                        google.removable -> OutlinedButton(onClick = vm::unlinkGoogle, enabled = !state.busy) { Text("Unlink Google") }
                        else -> Text("Google is your only way to sign in.")
                    }
                }
            }
            items(account.identitiesList.filter { it.provider != Provider.PROVIDER_GOOGLE }) {
                ListItem(headlineContent = { Text("Another sign-in method") })
            }
            item {
                HorizontalDivider()
                OutlinedButton(onClick = { open(Route.Sessions) }) { Text("Signed-in devices") }
                OutlinedButton(onClick = session::signOut) { Text("Sign out") }
                if (caps.deleteAccount) {
                    TextButton(onClick = { deleting = true }, enabled = !state.busy) {
                        Text("Delete account", color = MaterialTheme.colorScheme.error)
                    }
                }
            }
        }
    }

    renaming?.let { id ->
        val passkey = account.passkeysList.firstOrNull { it.id == id }
        var name by rememberSaveable(id) { mutableStateOf(passkey?.name.orEmpty()) }
        AlertDialog(
            onDismissRequest = { renaming = null },
            title = { Text("Rename passkey") },
            text = { OutlinedTextField(name, { name = it }, singleLine = true) },
            confirmButton = {
                TextButton(onClick = { passkey?.let { vm.rename(it, name.trim()) }; renaming = null }, enabled = name.isNotBlank()) {
                    Text("Save")
                }
            },
            dismissButton = { TextButton(onClick = { renaming = null }) { Text("Cancel") } },
        )
    }
    if (deleting) {
        AlertDialog(
            onDismissRequest = { deleting = false },
            title = { Text("Delete your account?") },
            text = { Text("This removes your account and signs out every device. It can't be undone.") },
            confirmButton = {
                TextButton(onClick = { deleting = false; vm.delete() }) { Text("Delete", color = MaterialTheme.colorScheme.error) }
            },
            dismissButton = { TextButton(onClick = { deleting = false }) { Text("Cancel") } },
        )
    }
}

private fun passkeyDetail(passkey: Passkey): String = buildString {
    append("Added ${date(passkey.createdAt)}")
    if (passkey.hasUsedAt()) append(" · last used ${date(passkey.usedAt)}")
    if (!passkey.usable) append(" · not usable here")
}

fun date(t: Timestamp): String =
    DateTimeFormatter.ofLocalizedDate(FormatStyle.MEDIUM)
        .withZone(ZoneId.systemDefault())
        .format(Instant.ofEpochSecond(t.seconds, t.nanos.toLong()))
