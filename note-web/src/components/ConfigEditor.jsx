import { useState } from 'react';

export default function ConfigEditor({ config, onSave }) {
  const [formData, setFormData] = useState(config);
  const [isEditing, setIsEditing] = useState(false);
  
  const handleInputChange = (field, value) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    onSave(formData);
    setIsEditing(false);
  };

  if (!config) {
    return null;
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="form-control">
        <label className="label">
          <span className="label-text">Editor</span>
        </label>
        <input
          type="text"
          value={formData.editor || ''}
          onChange={(e) => handleInputChange('editor', e.target.value)}
          className="input input-bordered"
          disabled={!isEditing}
        />
      </div>
      <div className="form-control">
        <label className="label">
          <span className="label-text">Date Format</span>
        </label>
        <input
          type="text"
          value={formData.date_format || ''}
          onChange={(e) => handleInputChange('date_format', e.target.value)}
          className="input input-bordered"
          disabled={!isEditing}
        />
      </div>

      {/* More fields as needed */}

      <div className="flex space-x-4 justify-end mt-6">
        {!isEditing && (
          <button
            type="button"
            className="btn btn-primary"
            onClick={() => setIsEditing(true)}
          >
            Edit
          </button>
        )}
        {isEditing && (
          <>
            <button
              type="button"
              className="btn btn-outline"
              onClick={() => {
                setFormData(config);
                setIsEditing(false);
              }}
            >
              Cancel
            </button>
            <button type="submit" className="btn btn-success">
              Save
            </button>
          </>
        )}
      </div>
    </form>
  );
}

