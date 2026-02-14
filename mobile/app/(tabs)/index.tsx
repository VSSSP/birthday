import { View, Text, StyleSheet } from "react-native";
import { useAuthStore } from "../../stores/authStore";

export default function HomeScreen() {
  const { user } = useAuthStore();

  return (
    <View style={styles.container}>
      <Text style={styles.greeting}>
        Hello, {user?.name || "there"}! 🎁
      </Text>
      <Text style={styles.subtitle}>
        Find the perfect birthday gift for your loved ones
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 24,
    backgroundColor: "#FFFFFF",
    justifyContent: "center",
    alignItems: "center",
  },
  greeting: {
    fontSize: 24,
    fontWeight: "bold",
    color: "#111827",
    marginBottom: 8,
    textAlign: "center",
  },
  subtitle: {
    fontSize: 16,
    color: "#6B7280",
    textAlign: "center",
  },
});
