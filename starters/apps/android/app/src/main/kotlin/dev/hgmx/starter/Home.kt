package dev.hgmx.starter

import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.lifecycle.compose.collectAsStateWithLifecycle

@Composable
fun HomeScreen(session: Session, open: (Route) -> Unit) {
    val gate by session.gate.collectAsStateWithLifecycle()
    val account = (gate as? Gate.SignedIn)?.account ?: return
    Centered {
        Text("Hi, ${account.displayName}", style = MaterialTheme.typography.headlineSmall)
        Button(onClick = { open(Route.Account) }) { Text("Account") }
    }
}
