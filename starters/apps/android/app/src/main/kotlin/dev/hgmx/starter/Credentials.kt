package dev.hgmx.starter

import android.content.Context
import androidx.credentials.CreatePublicKeyCredentialRequest
import androidx.credentials.CreatePublicKeyCredentialResponse
import androidx.credentials.CredentialManager
import androidx.credentials.CustomCredential
import androidx.credentials.GetCredentialRequest
import androidx.credentials.GetPublicKeyCredentialOption
import androidx.credentials.PublicKeyCredential
import androidx.credentials.exceptions.CreateCredentialCancellationException
import androidx.credentials.exceptions.CreateCredentialException
import androidx.credentials.exceptions.GetCredentialCancellationException
import androidx.credentials.exceptions.GetCredentialException
import androidx.credentials.exceptions.NoCredentialException
import com.google.android.libraries.identity.googleid.GetSignInWithGoogleOption
import com.google.android.libraries.identity.googleid.GoogleIdTokenCredential
import kotlinx.coroutines.CancellationException

object Credentials {
    suspend fun passkey(context: Context, options: String): String {
        val request = GetCredentialRequest(listOf(GetPublicKeyCredentialOption(options)))
        val credential = get(context, request)
        return (credential as? PublicKeyCredential)?.authenticationResponseJson
            ?: throw Problem("That sign-in option isn't a passkey.")
    }

    suspend fun createPasskey(context: Context, options: String): String {
        try {
            val response = CredentialManager.create(context)
                .createCredential(context, CreatePublicKeyCredentialRequest(options))
            return (response as CreatePublicKeyCredentialResponse).registrationResponseJson
        } catch (_: CreateCredentialCancellationException) {
            throw CancellationException("Passkey creation was cancelled.")
        } catch (e: CreateCredentialException) {
            throw Problem(e.errorMessage?.toString() ?: "Couldn't create a passkey.")
        }
    }

    suspend fun google(context: Context, clientId: String, nonce: String): String {
        val option = GetSignInWithGoogleOption.Builder(clientId).setNonce(nonce).build()
        val credential = get(context, GetCredentialRequest(listOf(option)))
        if (credential !is CustomCredential || credential.type != GoogleIdTokenCredential.TYPE_GOOGLE_ID_TOKEN_CREDENTIAL) {
            throw Problem("Google didn't return an account.")
        }
        return GoogleIdTokenCredential.createFrom(credential.data).idToken
    }

    private suspend fun get(context: Context, request: GetCredentialRequest) = try {
        CredentialManager.create(context).getCredential(context, request).credential
    } catch (_: GetCredentialCancellationException) {
        throw CancellationException("Sign-in was cancelled.")
    } catch (_: NoCredentialException) {
        throw Problem("No saved sign-in for this app on this device.")
    } catch (e: GetCredentialException) {
        throw Problem(e.errorMessage?.toString() ?: "Couldn't sign in.")
    }
}
