import { TouchableOpacity, Text, Alert, StyleSheet } from "react-native";
import {
  GoogleSignin,
  statusCodes,
} from "@react-native-google-signin/google-signin";
import { useAuthStore } from "../../stores/authStore";
import { GOOGLE_WEB_CLIENT_ID } from "../../constants/config";

GoogleSignin.configure({
  webClientId: GOOGLE_WEB_CLIENT_ID,
  offlineAccess: false,
});

export function GoogleSignInButton() {
  const { loginWithGoogle } = useAuthStore();

  const handleGoogleSignIn = async () => {
    try {
      await GoogleSignin.hasPlayServices();
      const response = await GoogleSignin.signIn();
      const idToken = response.data?.idToken;
      if (!idToken) {
        throw new Error("No ID token received from Google");
      }
      await loginWithGoogle(idToken);
    } catch (error: any) {
      if (error.code === statusCodes.SIGN_IN_CANCELLED) return;
      Alert.alert(
        "Google Sign-In Failed",
        error.message || "Something went wrong"
      );
    }
  };

  return (
    <TouchableOpacity onPress={handleGoogleSignIn} style={styles.button}>
      <Text style={styles.text}>Continue with Google</Text>
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  button: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "center",
    padding: 16,
    borderRadius: 12,
    borderWidth: 1,
    borderColor: "#E5E7EB",
    backgroundColor: "#FFFFFF",
  },
  text: {
    fontSize: 16,
    color: "#374151",
    fontWeight: "500",
  },
});
