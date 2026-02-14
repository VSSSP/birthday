import { useEffect, useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { useRecipientStore } from '../stores/recipientStore';
import { useAuthStore } from '../stores/authStore';
import RecipientCard from '../components/RecipientCard';
import Toast from '../components/Toast';
import Loading from '../components/Loading';
import styles from './Recipients.module.css';

export default function Recipients() {
  const { recipients, fetchRecipients, bulkDelete, isLoading } = useRecipientStore();
  const { user, logout } = useAuthStore();
  const navigate = useNavigate();
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [error, setError] = useState('');
  const [showProfile, setShowProfile] = useState(false);

  useEffect(() => {
    fetchRecipients();
  }, [fetchRecipients]);

  const toggleSelect = useCallback((id: string) => {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  }, []);

  const clearSelection = () => setSelectedIds(new Set());

  const handleBulkDelete = async () => {
    if (!confirm(`Delete ${selectedIds.size} recipient(s)?`)) return;
    try {
      await bulkDelete(Array.from(selectedIds));
      clearSelection();
    } catch {
      setError('Failed to delete recipients');
    }
  };

  const handleLogout = () => {
    if (confirm('Sign out?')) {
      logout();
      navigate('/login');
    }
  };

  const userInitial = user?.name?.charAt(0).toUpperCase() || 'U';

  return (
    <div className={styles.page}>
      {/* Header */}
      <header className={styles.header}>
        <h2>My Recipients</h2>
        <button className={styles.avatarBtn} onClick={() => setShowProfile(!showProfile)}>
          <span className={styles.avatarSmall}>{userInitial}</span>
        </button>
      </header>

      {/* Profile Dropdown */}
      {showProfile && (
        <div className={styles.profileDropdown}>
          <div className={styles.profileInfo}>
            <span className={styles.avatarLarge}>{userInitial}</span>
            <strong>{user?.name}</strong>
            <span className={styles.textMuted}>{user?.email}</span>
          </div>
          <button className={styles.logoutBtn} onClick={handleLogout}>
            Sign Out
          </button>
        </div>
      )}

      {/* Content */}
      <div className={styles.content}>
        {isLoading && recipients.length === 0 && <Loading />}

        {!isLoading && recipients.length === 0 && (
          <div className={styles.emptyState}>
            <div className={styles.emptyIcon}>&#127873;</div>
            <h3>No recipients yet</h3>
            <p>Add someone you want to find the perfect gift for</p>
          </div>
        )}

        {recipients.map((rec) => (
          <RecipientCard
            key={rec.id}
            recipient={rec}
            isSelected={selectedIds.has(rec.id)}
            isSelecting={selectedIds.size > 0}
            onSelect={() => toggleSelect(rec.id)}
            onClick={() => navigate(`/recipient/${rec.id}`)}
          />
        ))}
      </div>

      {/* Bulk Bar */}
      {selectedIds.size > 0 && (
        <div className={styles.bulkBar}>
          <span>{selectedIds.size} selected</span>
          <button className={styles.bulkDeleteBtn} onClick={handleBulkDelete}>
            Delete
          </button>
          <button className={styles.bulkCancelBtn} onClick={clearSelection}>
            Cancel
          </button>
        </div>
      )}

      {/* FAB */}
      <button className={styles.fab} onClick={() => navigate('/recipient/new')}>
        +
      </button>

      {error && <Toast message={error} type="error" onClose={() => setError('')} />}
    </div>
  );
}
