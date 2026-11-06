package com.example.fixture.notes

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.Button
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.ListItem
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.lifecycle.ViewModel
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.viewModelScope
import androidx.lifecycle.viewmodel.compose.viewModel
import fixture.api.v1.AddNoteRequest
import fixture.api.v1.ListNotesRequest
import fixture.api.v1.NoteServiceClient
import fixture.api.v1.Notes
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch

data class NotesState(
    val notes: Notes? = null,
    val draft: String = "",
    val busy: Boolean = false,
    val problem: Problem? = null,
)

class NotesViewModel(private val session: Session) : ViewModel() {
    private val notes = NoteServiceClient(session.api.client)
    private val _state = MutableStateFlow(NotesState())
    val state = _state.asStateFlow()

    init {
        act { session.call { notes.listNotes(ListNotesRequest.getDefaultInstance(), it) }.notes }
    }

    fun draft(text: String) = _state.update { it.copy(draft = text) }

    fun add() = act {
        val request = AddNoteRequest.newBuilder().setText(_state.value.draft.trim()).build()
        session.call { notes.addNote(request, it) }.notes.also { _state.update { s -> s.copy(draft = "") } }
    }

    private fun act(action: suspend () -> Notes) {
        if (_state.value.busy) return
        _state.update { it.copy(busy = true, problem = null) }
        viewModelScope.launch {
            try {
                val snapshot = action()
                _state.update { it.copy(notes = snapshot) }
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
fun HomeScreen(session: Session, open: (Route) -> Unit) {
    val vm = viewModel { NotesViewModel(session) }
    val state by vm.state.collectAsStateWithLifecycle()
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Notes") },
                actions = { TextButton(onClick = { open(Route.Account) }) { Text("Account") } },
            )
        },
    ) { padding ->
        LazyColumn(Modifier.padding(padding).padding(horizontal = 16.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
            item { ProblemText(state.problem, open) }
            val notes = state.notes
            if (notes == null) {
                item { if (state.busy) CircularProgressIndicator() }
                return@LazyColumn
            }
            if (notes.canAdd) {
                item {
                    Row(verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                        OutlinedTextField(state.draft, vm::draft, label = { Text("New note") }, modifier = Modifier.weight(1f))
                        Button(onClick = vm::add, enabled = !state.busy && state.draft.isNotBlank()) { Text("Add") }
                    }
                }
            }
            if (notes.notesCount == 0) item { Text("No notes yet.") }
            items(notes.notesList, key = { it.id }) { ListItem(headlineContent = { Text(it.text) }, modifier = Modifier.fillMaxWidth()) }
        }
    }
}
