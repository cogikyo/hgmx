package dev.hgmx.starter

import android.app.Application
import androidx.lifecycle.Lifecycle
import androidx.lifecycle.ProcessLifecycleOwner
import androidx.lifecycle.eventFlow
import kotlinx.coroutines.MainScope
import kotlinx.coroutines.flow.drop
import kotlinx.coroutines.flow.filter
import kotlinx.coroutines.flow.launchIn
import kotlinx.coroutines.flow.onEach
import kotlinx.coroutines.launch

class App : Application() {
    lateinit var session: Session
        private set

    override fun onCreate() {
        super.onCreate()
        val api = Api(BuildConfig.API_ORIGIN, BuildConfig.APP_ID, BuildConfig.VERSION_CODE)
        val scope = MainScope()
        session = Session(api, TokenStore(noBackupFilesDir), scope)
        scope.launch { session.start() }
        ProcessLifecycleOwner.get().lifecycle.eventFlow
            .filter { it == Lifecycle.Event.ON_START }
            .drop(1)
            // A stale-session cancellation must not stop future foreground checks.
            .onEach { scope.launch { session.recheck() } }
            .launchIn(scope)
    }
}
