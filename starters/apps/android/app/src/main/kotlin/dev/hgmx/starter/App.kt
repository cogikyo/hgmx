package dev.hgmx.starter

import android.app.Application
import kotlinx.coroutines.MainScope
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
    }
}
