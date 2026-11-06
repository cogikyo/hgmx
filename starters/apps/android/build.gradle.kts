plugins {
    alias(libs.plugins.android.application) apply false
    alias(libs.plugins.kotlin.compose) apply false
}

buildscript {
    dependencyLocking {
        lockAllConfigurations()
    }
}
