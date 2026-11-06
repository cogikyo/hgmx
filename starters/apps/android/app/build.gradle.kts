plugins {
    alias(libs.plugins.android.application)
    alias(libs.plugins.kotlin.compose)
}

android {
    namespace = "dev.hgmx.starter"
    compileSdk = 37
    buildToolsVersion = "36.0.0"

    defaultConfig {
        applicationId = "dev.hgmx.starter"
        minSdk = 28
        targetSdk = 37
        versionCode = 1
        versionName = "1.0"
        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"

        resValue("string", "app_name", "Starter")
        buildConfigField("String", "API_ORIGIN", "\"https://starter.hgmx.dev\"")
        buildConfigField("String", "APP_ID", "\"starter-android\"")
        buildConfigField("String", "GOOGLE_CLIENT_ID", "\"\"")
    }

    buildFeatures {
        buildConfig = true
        resValues = true
        compose = true
    }

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_21
        targetCompatibility = JavaVersion.VERSION_21
    }

    sourceSets.getByName("main") {
        java.directories += "src/main/generated"
        kotlin.directories += "src/main/generated"
    }
}

kotlin {
    jvmToolchain(21)
}

dependencyLocking {
    lockAllConfigurations()
}

dependencies {
    implementation(platform(libs.compose.bom))
    implementation(libs.compose.material3)
    implementation(libs.compose.ui.tooling.preview)
    implementation(libs.core.ktx)
    implementation(libs.activity.compose)
    implementation(libs.lifecycle.runtime.compose)
    implementation(libs.lifecycle.viewmodel.compose)
    implementation(libs.lifecycle.viewmodel.navigation3)
    implementation(libs.navigation3.runtime)
    implementation(libs.navigation3.ui)
    implementation(libs.credentials)
    implementation(libs.credentials.play.services)
    implementation(libs.googleid)
    implementation(libs.connect)
    implementation(libs.connect.okhttp)
    implementation(libs.connect.javalite)
    implementation(libs.protobuf.javalite)
    implementation(libs.okhttp)
    debugImplementation(libs.compose.ui.tooling)
    androidTestImplementation(libs.test.runner)
    androidTestImplementation(libs.test.junit)
}
