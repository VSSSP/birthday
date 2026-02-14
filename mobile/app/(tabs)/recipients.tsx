import { useEffect, useState, useCallback } from "react";
import {
  View,
  Text,
  FlatList,
  TouchableOpacity,
  StyleSheet,
  RefreshControl,
  Alert,
} from "react-native";
import { useRouter } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import { useRecipientStore } from "../../stores/recipientStore";
import { Recipient } from "../../types/recipient";

export default function RecipientsScreen() {
  const router = useRouter();
  const { recipients, isLoading, fetchRecipients, deleteRecipient, bulkDeleteRecipients } =
    useRecipientStore();
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [selectionMode, setSelectionMode] = useState(false);

  useEffect(() => {
    fetchRecipients();
  }, []);

  const onRefresh = useCallback(() => {
    fetchRecipients();
  }, []);

  const toggleSelection = (id: string) => {
    setSelectedIds((prev) =>
      prev.includes(id) ? prev.filter((i) => i !== id) : [...prev, id]
    );
  };

  const handleLongPress = (id: string) => {
    setSelectionMode(true);
    setSelectedIds([id]);
  };

  const handleBulkDelete = () => {
    Alert.alert(
      "Delete Recipients",
      `Are you sure you want to delete ${selectedIds.length} recipient(s)?`,
      [
        { text: "Cancel", style: "cancel" },
        {
          text: "Delete",
          style: "destructive",
          onPress: async () => {
            await bulkDeleteRecipients(selectedIds);
            setSelectedIds([]);
            setSelectionMode(false);
          },
        },
      ]
    );
  };

  const handleDelete = (id: string, name: string) => {
    Alert.alert("Delete", `Remove ${name} from your list?`, [
      { text: "Cancel", style: "cancel" },
      {
        text: "Delete",
        style: "destructive",
        onPress: () => deleteRecipient(id),
      },
    ]);
  };

  const cancelSelection = () => {
    setSelectionMode(false);
    setSelectedIds([]);
  };

  const renderItem = ({ item }: { item: Recipient }) => {
    const isSelected = selectedIds.includes(item.id);

    return (
      <TouchableOpacity
        style={[styles.card, isSelected && styles.cardSelected]}
        onPress={() => {
          if (selectionMode) {
            toggleSelection(item.id);
          } else {
            router.push(`/recipient/${item.id}`);
          }
        }}
        onLongPress={() => handleLongPress(item.id)}
      >
        <View style={styles.cardLeft}>
          <View style={[styles.avatarSmall, isSelected && styles.avatarSelected]}>
            {isSelected ? (
              <Ionicons name="checkmark" size={20} color="#FFFFFF" />
            ) : (
              <Text style={styles.avatarSmallText}>
                {item.name.charAt(0).toUpperCase()}
              </Text>
            )}
          </View>
          <View style={styles.cardInfo}>
            <Text style={styles.cardName}>{item.name}</Text>
            <Text style={styles.cardMeta}>
              {item.age} years old · {item.gender}
            </Text>
            <Text style={styles.cardBudget}>
              ${item.min_budget.toFixed(0)} - ${item.max_budget.toFixed(0)}
            </Text>
            {item.keywords.length > 0 && (
              <View style={styles.tagsRow}>
                {item.keywords.slice(0, 3).map((kw, i) => (
                  <View key={i} style={styles.tag}>
                    <Text style={styles.tagText}>{kw}</Text>
                  </View>
                ))}
                {item.keywords.length > 3 && (
                  <Text style={styles.moreTag}>
                    +{item.keywords.length - 3}
                  </Text>
                )}
              </View>
            )}
          </View>
        </View>

        {!selectionMode && (
          <TouchableOpacity
            onPress={() => handleDelete(item.id, item.name)}
            style={styles.deleteButton}
          >
            <Ionicons name="trash-outline" size={20} color="#EF4444" />
          </TouchableOpacity>
        )}
      </TouchableOpacity>
    );
  };

  return (
    <View style={styles.container}>
      {selectionMode && (
        <View style={styles.selectionBar}>
          <TouchableOpacity onPress={cancelSelection}>
            <Text style={styles.cancelText}>Cancel</Text>
          </TouchableOpacity>
          <Text style={styles.selectedCount}>
            {selectedIds.length} selected
          </Text>
          <TouchableOpacity onPress={handleBulkDelete}>
            <Ionicons name="trash" size={24} color="#EF4444" />
          </TouchableOpacity>
        </View>
      )}

      <FlatList
        data={recipients}
        keyExtractor={(item) => item.id}
        renderItem={renderItem}
        refreshControl={
          <RefreshControl refreshing={isLoading} onRefresh={onRefresh} />
        }
        contentContainerStyle={[
          styles.list,
          recipients.length === 0 && styles.emptyList,
        ]}
        ListEmptyComponent={
          <View style={styles.empty}>
            <Text style={styles.emptyEmoji}>👤</Text>
            <Text style={styles.emptyTitle}>No recipients yet</Text>
            <Text style={styles.emptySubtitle}>
              Add someone you want to buy a gift for
            </Text>
          </View>
        }
      />

      <TouchableOpacity
        style={styles.fab}
        onPress={() => router.push("/recipient/new")}
      >
        <Ionicons name="add" size={28} color="#FFFFFF" />
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#F9FAFB",
  },
  selectionBar: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    paddingHorizontal: 20,
    paddingVertical: 12,
    backgroundColor: "#FFFFFF",
    borderBottomWidth: 1,
    borderBottomColor: "#E5E7EB",
  },
  cancelText: {
    color: "#7C3AED",
    fontSize: 16,
    fontWeight: "500",
  },
  selectedCount: {
    fontSize: 16,
    fontWeight: "600",
    color: "#111827",
  },
  list: {
    padding: 16,
    paddingBottom: 100,
  },
  emptyList: {
    flex: 1,
    justifyContent: "center",
  },
  card: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    backgroundColor: "#FFFFFF",
    borderRadius: 16,
    padding: 16,
    marginBottom: 12,
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.05,
    shadowRadius: 4,
    elevation: 1,
  },
  cardSelected: {
    borderWidth: 2,
    borderColor: "#7C3AED",
  },
  cardLeft: {
    flexDirection: "row",
    alignItems: "center",
    flex: 1,
  },
  avatarSmall: {
    width: 48,
    height: 48,
    borderRadius: 24,
    backgroundColor: "#EDE9FE",
    justifyContent: "center",
    alignItems: "center",
    marginRight: 14,
  },
  avatarSelected: {
    backgroundColor: "#7C3AED",
  },
  avatarSmallText: {
    fontSize: 18,
    fontWeight: "bold",
    color: "#7C3AED",
  },
  cardInfo: {
    flex: 1,
  },
  cardName: {
    fontSize: 16,
    fontWeight: "600",
    color: "#111827",
    marginBottom: 2,
  },
  cardMeta: {
    fontSize: 13,
    color: "#6B7280",
    marginBottom: 2,
  },
  cardBudget: {
    fontSize: 13,
    color: "#059669",
    fontWeight: "500",
    marginBottom: 6,
  },
  tagsRow: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: 6,
  },
  tag: {
    backgroundColor: "#F3F4F6",
    paddingHorizontal: 8,
    paddingVertical: 3,
    borderRadius: 6,
  },
  tagText: {
    fontSize: 11,
    color: "#6B7280",
  },
  moreTag: {
    fontSize: 11,
    color: "#9CA3AF",
    alignSelf: "center",
  },
  deleteButton: {
    padding: 8,
  },
  empty: {
    alignItems: "center",
  },
  emptyEmoji: {
    fontSize: 48,
    marginBottom: 16,
  },
  emptyTitle: {
    fontSize: 18,
    fontWeight: "600",
    color: "#111827",
    marginBottom: 4,
  },
  emptySubtitle: {
    fontSize: 14,
    color: "#6B7280",
  },
  fab: {
    position: "absolute",
    right: 20,
    bottom: 24,
    width: 56,
    height: 56,
    borderRadius: 28,
    backgroundColor: "#7C3AED",
    justifyContent: "center",
    alignItems: "center",
    shadowColor: "#7C3AED",
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.3,
    shadowRadius: 8,
    elevation: 6,
  },
});
