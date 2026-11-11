package dev.hgmx.starter

import android.content.Intent
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.CompositionLocalProvider
import androidx.compose.runtime.getValue
import androidx.compose.runtime.key
import androidx.compose.runtime.mutableStateListOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp
import androidx.core.net.toUri
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.viewmodel.compose.LocalViewModelStoreOwner
import androidx.lifecycle.viewmodel.compose.rememberViewModelStoreOwner
import androidx.lifecycle.viewmodel.navigation3.rememberViewModelStoreNavEntryDecorator
import androidx.navigation3.runtime.NavKey
import androidx.navigation3.runtime.entryProvider
import androidx.navigation3.runtime.rememberSaveableStateHolderNavEntryDecorator
import androidx.navigation3.ui.NavDisplay
import kotlinx.coroutines.launch

sealed interface Route : NavKey {
    data object Home : Route
    data object Account : Route
    data object Sessions : Route
    data object Reauth : Route
}

@Composable
fun Root(session: Session) {
    val gate by session.gate.collectAsStateWithLifecycle()
    Surface(Modifier.fillMaxSize()) {
        when (val g = gate) {
            Gate.Starting -> Centered { CircularProgressIndicator() }
            is Gate.Unreachable -> {
                val scope = rememberCoroutineScope()
                Centered {
                    Text(g.problem.message.orEmpty())
                    Button(onClick = { scope.launch { session.start() } }) { Text("Try again") }
                }
            }
            is Gate.Update -> UpdateScreen(g)
            is Gate.SignedOut -> Scoped { SignInScreen(session, notice = g.notice, onDone = {}) }
            is Gate.SignedIn -> key(g.account.id) { Scoped { SignedIn(session) } }
        }
    }
}

@Composable
private fun Scoped(content: @Composable () -> Unit) {
    CompositionLocalProvider(LocalViewModelStoreOwner provides rememberViewModelStoreOwner(), content = content)
}

@Composable
private fun UpdateScreen(update: Gate.Update) {
    val context = LocalContext.current
    val store = update.storeUrl.toUri().takeIf { it.scheme == "https" }
    Centered {
        Text("Update required", style = MaterialTheme.typography.headlineSmall)
        Text(update.message)
        if (store != null) {
            Button(onClick = { context.startActivity(Intent(Intent.ACTION_VIEW, store)) }) { Text("Open the store") }
        }
    }
}

@Composable
private fun SignedIn(session: Session) {
    val stack = remember { mutableStateListOf<Route>(Route.Home) }
    val open: (Route) -> Unit = { stack.add(it) }
    val back: () -> Unit = { if (stack.size > 1) stack.removeAt(stack.lastIndex) }
    NavDisplay(
        backStack = stack,
        onBack = back,
        entryDecorators = listOf(
            rememberSaveableStateHolderNavEntryDecorator(),
            rememberViewModelStoreNavEntryDecorator(),
        ),
        entryProvider = entryProvider {
            entry<Route.Home> { HomeScreen(session, open) }
            entry<Route.Account> { AccountScreen(session, open, back) }
            entry<Route.Sessions> { SessionsScreen(session, open, back) }
            entry<Route.Reauth> { SignInScreen(session, notice = "Sign in again to continue.", onDone = back) }
        },
    )
}

@Composable
fun Centered(content: @Composable () -> Unit) {
    Column(
        Modifier.fillMaxSize().padding(24.dp),
        verticalArrangement = Arrangement.spacedBy(16.dp, Alignment.CenterVertically),
        horizontalAlignment = Alignment.CenterHorizontally,
    ) {
        content()
    }
}

@Composable
fun ProblemText(problem: Problem?, open: (Route) -> Unit) {
    if (problem == null) return
    Text(problem.message.orEmpty(), color = MaterialTheme.colorScheme.error)
    if (problem.reauth) {
        Button(onClick = { open(Route.Reauth) }) { Text("Sign in again") }
    }
}
