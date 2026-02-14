import { useEffect, useState } from "react";
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  Alert,
  ScrollView,
  KeyboardAvoidingView,
  Platform,
  ActivityIndicator,
} from "react-native";
import { useRouter, useLocalSearchParams, Stack } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import { useRecipientStore } from "../../stores/recipientStore";
import { recipientService } from "../../services/recipientService";
import { Recipient } from "../../types/recipient";

const GENDER_OPTIONS = ["male", "female", "other"];

const SUGGESTED_KEYWORDS = [
  "nerd",
  "geek",
  "sports",
  "music",
  "cooking",
  "reading",
  "gaming",
  "travel",
  "fitness",
  "tech",
  "fashion",
  "art",
  "movies",
  "outdoor",
  "pets",
];

export default function RecipientDetailScreen() {
  const router = useRouter();
  const { id } = useLocalSearchParams<{ id: string }>();
  const { updateRecipient, deleteRecipient } = useRecipientStore();

  const [recipient, setRecipient] = useState<Recipient | null>(null);
  const [name, setName] = useState("");
  const [age, setAge] = useState("");
  const [gender, setGender] = useState("other");
  const [minBudget, setMinBudget] = useState("");
  const [maxBudget, setMaxBudget] = useState("");
  const [keywords, setKeywords] = useState<string[]>([]);
  const [customKeyword, setCustomKeyword] = useState("");
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    loadRecipient();
  }, [id]);

  const loadRecipient = async () => {
    try {
      const data = await recipientService.getById(id!);
      setRecipient(data);
      setName(data.name);
      setAge(data.age.toString());
      setGender(data.gender);
      setMinBudget(data.min_budget.toString());
      setMaxBudget(data.max_budget.toString());
      setKeywords(data.keywords || []);
    } catch {
      Alert.alert("Error", "Failed to load recipient");
      router.back();
    } finally {
      setLoading(false);
    }
  };

  const addKeyword = (kw: string) => {
    const normalized = kw.trim().toLowerCase();
    if (normalized && !keywords.includes(normalized)) {
      setKeywords([...keywords, normalized]);
    }
    setCustomKeyword("");
  };

  const removeKeyword = (kw: string) => {
    setKeywords(keywords.filter((k) => k !== kw));
  };

  const handleSave = async () => {
    if (!name.trim()) {
      Alert.alert("Error", "Name is required");
      return;
    }
    setSaving(true);
    try {
      await updateRecipient(id!, {
        name: name.trim(),
        age: parseInt(age) || 0,
        gender,
        min_budget: parseFloat(minBudget) || 0,
        max_budget: parseFloat(maxBudget) || 0,
        keywords,
      });
      router.back();
    } catch (error: any) {
      Alert.alert("Error", error.response?.data?.error || "Failed to update");
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = () => {
    Alert.alert("Delete", `Remove ${name} from your list?`, [
      { text: "Cancel", style: "cancel" },
      {
        text: "Delete",
        style: "destructive",
        onPress: async () => {
          await deleteRecipient(id!);
          router.back();
        },
      },
    ]);
  };

  if (loading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color="#7C3AED" />
      </View>
    );
  }

  return (
    <>
      <Stack.Screen
        options={{
          title: "Edit Recipient",
          headerStyle: { backgroundColor: "#FFFFFF" },
          headerTitleStyle: { fontWeight: "600" },
          headerRight: () => (
            <TouchableOpacity onPress={handleDelete} style={{ padding: 8 }}>
              <Ionicons name="trash-outline" size={22} color="#EF4444" />
            </TouchableOpacity>
          ),
        }}
      />
      <KeyboardAvoidingView
        style={styles.container}
        behavior={Platform.OS === "ios" ? "padding" : "height"}
      >
        <ScrollView
          contentContainerStyle={styles.scroll}
          keyboardShouldPersistTaps="handled"
        >
          <View style={styles.section}>
            <Text style={styles.label}>Name *</Text>
            <TextInput
              value={name}
              onChangeText={setName}
              style={styles.input}
              placeholderTextColor="#9CA3AF"
            />
          </View>

          <View style={styles.row}>
            <View style={[styles.section, { flex: 1, marginRight: 8 }]}>
              <Text style={styles.label}>Age</Text>
              <TextInput
                value={age}
                onChangeText={setAge}
                keyboardType="numeric"
                style={styles.input}
                placeholderTextColor="#9CA3AF"
              />
            </View>
            <View style={[styles.section, { flex: 1, marginLeft: 8 }]}>
              <Text style={styles.label}>Gender</Text>
              <View style={styles.genderRow}>
                {GENDER_OPTIONS.map((g) => (
                  <TouchableOpacity
                    key={g}
                    style={[
                      styles.genderChip,
                      gender === g && styles.genderChipActive,
                    ]}
                    onPress={() => setGender(g)}
                  >
                    <Text
                      style={[
                        styles.genderChipText,
                        gender === g && styles.genderChipTextActive,
                      ]}
                    >
                      {g.charAt(0).toUpperCase() + g.slice(1)}
                    </Text>
                  </TouchableOpacity>
                ))}
              </View>
            </View>
          </View>

          <View style={styles.section}>
            <Text style={styles.label}>Budget Range ($)</Text>
            <View style={styles.row}>
              <TextInput
                value={minBudget}
                onChangeText={setMinBudget}
                keyboardType="decimal-pad"
                style={[styles.input, { flex: 1, marginRight: 8 }]}
                placeholderTextColor="#9CA3AF"
              />
              <Text style={styles.budgetSeparator}>-</Text>
              <TextInput
                value={maxBudget}
                onChangeText={setMaxBudget}
                keyboardType="decimal-pad"
                style={[styles.input, { flex: 1, marginLeft: 8 }]}
                placeholderTextColor="#9CA3AF"
              />
            </View>
          </View>

          <View style={styles.section}>
            <Text style={styles.label}>Keywords / Interests</Text>

            {keywords.length > 0 && (
              <View style={styles.selectedTags}>
                {keywords.map((kw) => (
                  <TouchableOpacity
                    key={kw}
                    style={styles.selectedTag}
                    onPress={() => removeKeyword(kw)}
                  >
                    <Text style={styles.selectedTagText}>{kw}</Text>
                    <Ionicons name="close" size={14} color="#7C3AED" />
                  </TouchableOpacity>
                ))}
              </View>
            )}

            <View style={styles.customKeywordRow}>
              <TextInput
                placeholder="Add custom keyword..."
                value={customKeyword}
                onChangeText={setCustomKeyword}
                onSubmitEditing={() => addKeyword(customKeyword)}
                returnKeyType="done"
                style={[styles.input, { flex: 1, marginBottom: 0 }]}
                placeholderTextColor="#9CA3AF"
              />
              <TouchableOpacity
                style={styles.addButton}
                onPress={() => addKeyword(customKeyword)}
              >
                <Ionicons name="add" size={24} color="#7C3AED" />
              </TouchableOpacity>
            </View>

            <View style={styles.suggestedTags}>
              {SUGGESTED_KEYWORDS.filter((kw) => !keywords.includes(kw)).map(
                (kw) => (
                  <TouchableOpacity
                    key={kw}
                    style={styles.suggestedTag}
                    onPress={() => addKeyword(kw)}
                  >
                    <Text style={styles.suggestedTagText}>{kw}</Text>
                    <Ionicons name="add" size={14} color="#6B7280" />
                  </TouchableOpacity>
                )
              )}
            </View>
          </View>

          <TouchableOpacity
            onPress={handleSave}
            disabled={saving}
            style={[styles.saveButton, saving && styles.saveButtonDisabled]}
          >
            <Text style={styles.saveButtonText}>
              {saving ? "Saving..." : "Save Changes"}
            </Text>
          </TouchableOpacity>
        </ScrollView>
      </KeyboardAvoidingView>
    </>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#F9FAFB",
  },
  loadingContainer: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    backgroundColor: "#F9FAFB",
  },
  scroll: {
    padding: 20,
    paddingBottom: 40,
  },
  section: {
    marginBottom: 20,
  },
  label: {
    fontSize: 14,
    fontWeight: "600",
    color: "#374151",
    marginBottom: 8,
  },
  input: {
    borderWidth: 1,
    borderColor: "#E5E7EB",
    backgroundColor: "#FFFFFF",
    padding: 14,
    borderRadius: 12,
    fontSize: 16,
    color: "#111827",
  },
  row: {
    flexDirection: "row",
    alignItems: "center",
  },
  budgetSeparator: {
    fontSize: 18,
    color: "#9CA3AF",
  },
  genderRow: {
    flexDirection: "row",
    gap: 6,
  },
  genderChip: {
    paddingHorizontal: 12,
    paddingVertical: 8,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: "#E5E7EB",
    backgroundColor: "#FFFFFF",
  },
  genderChipActive: {
    backgroundColor: "#7C3AED",
    borderColor: "#7C3AED",
  },
  genderChipText: {
    fontSize: 13,
    color: "#6B7280",
  },
  genderChipTextActive: {
    color: "#FFFFFF",
    fontWeight: "600",
  },
  selectedTags: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: 8,
    marginBottom: 12,
  },
  selectedTag: {
    flexDirection: "row",
    alignItems: "center",
    backgroundColor: "#EDE9FE",
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 16,
    gap: 4,
  },
  selectedTagText: {
    fontSize: 13,
    color: "#7C3AED",
    fontWeight: "500",
  },
  customKeywordRow: {
    flexDirection: "row",
    alignItems: "center",
    gap: 8,
    marginBottom: 16,
  },
  addButton: {
    width: 48,
    height: 48,
    borderRadius: 12,
    borderWidth: 1,
    borderColor: "#E5E7EB",
    backgroundColor: "#FFFFFF",
    justifyContent: "center",
    alignItems: "center",
  },
  suggestedTags: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: 8,
  },
  suggestedTag: {
    flexDirection: "row",
    alignItems: "center",
    backgroundColor: "#F3F4F6",
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 16,
    gap: 4,
  },
  suggestedTagText: {
    fontSize: 13,
    color: "#6B7280",
  },
  saveButton: {
    backgroundColor: "#7C3AED",
    padding: 16,
    borderRadius: 12,
    alignItems: "center",
    marginTop: 8,
  },
  saveButtonDisabled: {
    opacity: 0.6,
  },
  saveButtonText: {
    color: "#FFFFFF",
    fontSize: 16,
    fontWeight: "600",
  },
});
