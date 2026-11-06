package dev.hgmx.starter

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.ListItem
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.lifecycle.ViewModel
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.viewModelScope
import androidx.lifecycle.viewmodel.compose.viewModel
import hgmx.auth.v1.ClientKind
import hgmx.auth.v1.ListSessionsRequest
import hgmx.auth.v1.RevokeOtherSessionsRequest
import hgmx.auth.v1.RevokeSessionRequest
import hgmx.auth.v1.Session as Device
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch

data class SessionsState(
    val devices: List<Device>? = null,
    val busy: Boolean = false,
    val problem: Problem? = null,
)

class SessionsViewModel(private val session: Session) : ViewModel() {
    private val api = session.api
    private val _state = MutableStateFlow(SessionsState())
    val state = _state.asStateFlow()

    init {
        act { devices() }
    }

    fun revoke(device: Device) = act {
        session.call { api.auth.revokeSession(RevokeSessionRequest.newBuilder().setSessionId(device.id).build(), it) }
        devices()
    }

    fun revokeOthers() = act {
        session.call { api.auth.revokeOtherSessions(RevokeOtherSessionsRequest.getDefaultInstance(), it) }
        devices()
    }

    private suspend fun devices() =
        session.call { api.auth.listSessions(ListSessionsRequest.getDefaultInstance(), it) }.sessionsList

    private fun act(action: suspend () -> List<Device>) {
        if (_state.value.busy) return
        _state.update { it.copy(busy = true, problem = null) }
        viewModelScope.launch {
            try {
                val devices = action()
                _state.update { it.copy(devices = devices) }
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
fun SessionsScreen(session: Session, open: (Route) -> Unit, back: () -> Unit) {
    val vm = viewModel { SessionsViewModel(session) }
    val state by vm.state.collectAsStateWithLifecycle()
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Signed-in devices") },
                navigationIcon = { TextButton(onClick = back) { Text("Back") } },
            )
        },
    ) { padding ->
        LazyColumn(Modifier.padding(padding).padding(horizontal = 16.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
            item { ProblemText(state.problem, open) }
            val devices = state.devices
            if (devices == null) {
                item { if (state.busy) CircularProgressIndicator() }
                return@LazyColumn
            }
            items(devices, key = { it.id }) { device ->
                ListItem(
                    headlineContent = { Text(device.label.ifEmpty { kind(device.kind) }) },
                    supportingContent = { Text("${kind(device.kind)} · last seen ${date(device.seenAt)}") },
                    trailingContent = {
                        if (device.current) {
                            Text("This device")
                        } else {
                            TextButton(onClick = { vm.revoke(device) }, enabled = !state.busy) { Text("Sign out") }
                        }
                    },
                )
            }
            if (devices.count { !it.current } > 0) {
                item {
                    OutlinedButton(onClick = vm::revokeOthers, enabled = !state.busy) { Text("Sign out everywhere else") }
                }
            }
        }
    }
}

private fun kind(kind: ClientKind) = when (kind) {
    ClientKind.CLIENT_KIND_BROWSER -> "Browser"
    ClientKind.CLIENT_KIND_IOS -> "iPhone or iPad"
    ClientKind.CLIENT_KIND_ANDROID -> "Android"
    else -> "Other device"
}
