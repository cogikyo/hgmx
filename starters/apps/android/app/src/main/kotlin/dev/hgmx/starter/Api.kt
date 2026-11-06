package dev.hgmx.starter

import androidx.core.net.toUri
import com.connectrpc.Headers
import com.connectrpc.ProtocolClientConfig
import com.connectrpc.extensions.GoogleJavaLiteProtobufStrategy
import com.connectrpc.impl.ProtocolClient
import com.connectrpc.okhttp.ConnectOkHttpClient
import hgmx.app.v1.AppServiceClient
import hgmx.auth.v1.AuthServiceClient
import okhttp3.OkHttpClient

class Api(origin: String, private val appId: String, private val build: Int) {
    val origin = exactOrigin(origin)

    val client = ProtocolClient(
        httpClient = ConnectOkHttpClient(
            ConnectOkHttpClient.configureClient(OkHttpClient.Builder())
                .followRedirects(false)
                .followSslRedirects(false)
                .build(),
        ),
        config = ProtocolClientConfig(
            host = this.origin,
            serializationStrategy = GoogleJavaLiteProtobufStrategy(),
        ),
    )
    val app = AppServiceClient(client)
    val auth = AuthServiceClient(client)

    fun headers(token: String?): Headers = buildMap {
        put("HGMX-App-ID", listOf(appId))
        put("HGMX-App-Build", listOf(build.toString()))
        if (token != null) put("Authorization", listOf("Bearer $token"))
    }
}

private fun exactOrigin(origin: String): String {
    val uri = origin.toUri()
    val host = uri.host
    val local = BuildConfig.DEBUG && uri.scheme == "http" && host == "localhost"
    require(
        (uri.scheme == "https" || local) && !host.isNullOrEmpty() && uri.userInfo == null &&
            uri.path.isNullOrEmpty() && uri.query == null && uri.fragment == null,
    ) { "API_ORIGIN must be an exact https origin, such as https://example.com" }
    return origin
}
