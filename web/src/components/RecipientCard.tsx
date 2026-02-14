import type { Recipient } from '../types/recipient';
import styles from './RecipientCard.module.css';

interface Props {
  recipient: Recipient;
  isSelected: boolean;
  isSelecting: boolean;
  onSelect: () => void;
  onClick: () => void;
}

export default function RecipientCard({
  recipient,
  isSelected,
  isSelecting,
  onSelect,
  onClick,
}: Props) {
  const handleClick = () => {
    if (isSelecting) {
      onSelect();
    } else {
      onClick();
    }
  };

  const handleLongPress = (e: React.MouseEvent | React.TouchEvent) => {
    e.preventDefault();
    onSelect();
  };

  return (
    <div
      className={`${styles.card} ${isSelected ? styles.selected : ''}`}
      onClick={handleClick}
      onContextMenu={handleLongPress}
    >
      <div className={styles.cardHeader}>
        <div>
          <div className={styles.name}>{recipient.name}</div>
          <div className={styles.meta}>
            {recipient.age && `${recipient.age} years`}
            {recipient.age && recipient.gender && ' · '}
            {recipient.gender}
          </div>
        </div>
        {isSelecting && (
          <div className={`${styles.checkbox} ${isSelected ? styles.checked : ''}`}>
            {isSelected && '✓'}
          </div>
        )}
      </div>

      {(recipient.min_budget > 0 || recipient.max_budget > 0) && (
        <div className={styles.budget}>
          ${recipient.min_budget?.toFixed(0) || '0'} — ${recipient.max_budget?.toFixed(0) || '0'}
        </div>
      )}

      {recipient.keywords && recipient.keywords.length > 0 && (
        <div className={styles.keywords}>
          {recipient.keywords.map((kw) => (
            <span key={kw} className={styles.keyword}>
              {kw}
            </span>
          ))}
        </div>
      )}
    </div>
  );
}
