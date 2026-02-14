import { useMemo, useState, type FormEvent } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useRecipientStore } from '../stores/recipientStore';
import Toast from '../components/Toast';
import Loading from '../components/Loading';
import styles from './RecipientForm.module.css';

const SUGGESTED_KEYWORDS = [
  'nerd', 'geek', 'sports', 'music', 'cooking', 'reading',
  'gaming', 'travel', 'fitness', 'tech', 'fashion', 'art',
  'movies', 'outdoor', 'pets',
];

export default function RecipientForm() {
  const { id } = useParams<{ id: string }>();
  const isEditing = !!id && id !== 'new';
  const navigate = useNavigate();
  const { recipients, createRecipient, updateRecipient, deleteRecipient, isLoading } =
    useRecipientStore();

  const existing = useMemo(
    () => (isEditing ? recipients.find((r) => r.id === id) : undefined),
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [id]
  );

  const [name, setName] = useState(existing?.name ?? '');
  const [age, setAge] = useState(existing?.age ? String(existing.age) : '');
  const [gender, setGender] = useState(existing?.gender ?? '');
  const [minBudget, setMinBudget] = useState(existing?.min_budget ? String(existing.min_budget) : '');
  const [maxBudget, setMaxBudget] = useState(existing?.max_budget ? String(existing.max_budget) : '');
  const [keywords, setKeywords] = useState<string[]>(existing?.keywords ?? []);
  const [keywordInput, setKeywordInput] = useState('');
  const [error, setError] = useState('');

  const addKeyword = (kw?: string) => {
    const value = (kw || keywordInput).trim().toLowerCase();
    if (value && !keywords.includes(value)) {
      setKeywords([...keywords, value]);
    }
    setKeywordInput('');
  };

  const removeKeyword = (kw: string) => {
    setKeywords(keywords.filter((k) => k !== kw));
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError('');

    if (!name.trim()) {
      setError('Name is required');
      return;
    }

    const data = {
      name: name.trim(),
      age: age ? parseInt(age) : 0,
      gender: gender || 'other',
      min_budget: minBudget ? parseFloat(minBudget) : 0,
      max_budget: maxBudget ? parseFloat(maxBudget) : 0,
      keywords,
    };

    try {
      if (isEditing) {
        await updateRecipient(id!, data);
      } else {
        await createRecipient(data);
      }
      navigate('/');
    } catch {
      setError(isEditing ? 'Failed to update recipient' : 'Failed to create recipient');
    }
  };

  const handleDelete = async () => {
    if (!confirm(`Delete ${name}?`)) return;
    try {
      await deleteRecipient(id!);
      navigate('/');
    } catch {
      setError('Failed to delete recipient');
    }
  };

  return (
    <div className={styles.page}>
      <header className={styles.header}>
        <button className={styles.backBtn} onClick={() => navigate('/')}>
          &larr;
        </button>
        <h2>{isEditing ? 'Edit Recipient' : 'New Recipient'}</h2>
        {isEditing && (
          <button className={styles.deleteBtn} onClick={handleDelete} title="Delete">
            &#128465;
          </button>
        )}
        {!isEditing && <div style={{ width: 32 }} />}
      </header>

      <form className={styles.content} onSubmit={handleSubmit}>
        {isLoading && <Loading />}

        <div className={styles.formGroup}>
          <label htmlFor="rec-name">Name *</label>
          <input
            id="rec-name"
            type="text"
            placeholder="Person's name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        </div>

        <div className={styles.row}>
          <div className={`${styles.formGroup} ${styles.flex1}`}>
            <label htmlFor="rec-age">Age</label>
            <input
              id="rec-age"
              type="number"
              placeholder="Age"
              min="0"
              max="150"
              value={age}
              onChange={(e) => setAge(e.target.value)}
            />
          </div>
          <div className={`${styles.formGroup} ${styles.flex1}`}>
            <label>Gender</label>
            <div className={styles.genderChips}>
              {['male', 'female', 'other'].map((g) => (
                <button
                  key={g}
                  type="button"
                  className={`${styles.chip} ${gender === g ? styles.chipActive : ''}`}
                  onClick={() => setGender(g)}
                >
                  {g.charAt(0).toUpperCase() + g.slice(1)}
                </button>
              ))}
            </div>
          </div>
        </div>

        <div className={styles.formGroup}>
          <label>Budget Range</label>
          <div className={styles.row}>
            <div className={`${styles.inputPrefix} ${styles.flex1}`}>
              <span>$</span>
              <input
                type="number"
                placeholder="Min"
                min="0"
                step="0.01"
                value={minBudget}
                onChange={(e) => setMinBudget(e.target.value)}
              />
            </div>
            <span className={styles.budgetSep}>to</span>
            <div className={`${styles.inputPrefix} ${styles.flex1}`}>
              <span>$</span>
              <input
                type="number"
                placeholder="Max"
                min="0"
                step="0.01"
                value={maxBudget}
                onChange={(e) => setMaxBudget(e.target.value)}
              />
            </div>
          </div>
        </div>

        <div className={styles.formGroup}>
          <label>Keywords / Interests</label>

          {keywords.length > 0 && (
            <div className={styles.keywordsContainer}>
              {keywords.map((kw) => (
                <span key={kw} className={styles.keywordTag}>
                  {kw}
                  <button type="button" onClick={() => removeKeyword(kw)}>
                    &times;
                  </button>
                </span>
              ))}
            </div>
          )}

          <div className={styles.keywordInputRow}>
            <input
              type="text"
              placeholder="Add a keyword..."
              value={keywordInput}
              onChange={(e) => setKeywordInput(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === 'Enter') {
                  e.preventDefault();
                  addKeyword();
                }
              }}
            />
            <button type="button" className={styles.addBtn} onClick={() => addKeyword()}>
              Add
            </button>
          </div>

          <div className={styles.suggested}>
            <span className={styles.suggestedLabel}>Suggestions:</span>
            <div className={styles.suggestedList}>
              {SUGGESTED_KEYWORDS.filter((kw) => !keywords.includes(kw)).map((kw) => (
                <button
                  key={kw}
                  type="button"
                  className={styles.chipSuggested}
                  onClick={() => addKeyword(kw)}
                >
                  {kw}
                </button>
              ))}
            </div>
          </div>
        </div>

        <button type="submit" className={styles.saveBtn} disabled={isLoading}>
          {isLoading ? 'Saving...' : 'Save Recipient'}
        </button>
      </form>

      {error && <Toast message={error} type="error" onClose={() => setError('')} />}
    </div>
  );
}
