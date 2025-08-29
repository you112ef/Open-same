-- Open-Same Database Initialization Script
-- This script creates the database schema and initial data

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    avatar TEXT,
    bio TEXT,
    is_admin BOOLEAN DEFAULT FALSE,
    is_verified BOOLEAN DEFAULT FALSE,
    is_banned BOOLEAN DEFAULT FALSE,
    last_login TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create content table
CREATE TABLE IF NOT EXISTS content (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    content TEXT,
    metadata JSONB,
    is_public BOOLEAN DEFAULT FALSE,
    is_template BOOLEAN DEFAULT FALSE,
    version INTEGER DEFAULT 1,
    creator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create content_versions table
CREATE TABLE IF NOT EXISTS content_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    title VARCHAR(255),
    content TEXT,
    metadata JSONB,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create collaborations table
CREATE TABLE IF NOT EXISTS collaborations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL CHECK (role IN ('owner', 'editor', 'viewer')),
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'declined')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, content_id)
);

-- Create shared_content table
CREATE TABLE IF NOT EXISTS shared_content (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    shared_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    share_type VARCHAR(50) NOT NULL CHECK (share_type IN ('public', 'link', 'embed')),
    share_url VARCHAR(500) UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE,
    view_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

CREATE INDEX IF NOT EXISTS idx_content_creator_id ON content(creator_id);
CREATE INDEX IF NOT EXISTS idx_content_type ON content(type);
CREATE INDEX IF NOT EXISTS idx_content_is_public ON content(is_public);
CREATE INDEX IF NOT EXISTS idx_content_created_at ON content(created_at);
CREATE INDEX IF NOT EXISTS idx_content_metadata ON content USING GIN(metadata);

CREATE INDEX IF NOT EXISTS idx_content_versions_content_id ON content_versions(content_id);
CREATE INDEX IF NOT EXISTS idx_content_versions_version ON content_versions(version);

CREATE INDEX IF NOT EXISTS idx_collaborations_user_id ON collaborations(user_id);
CREATE INDEX IF NOT EXISTS idx_collaborations_content_id ON collaborations(content_id);
CREATE INDEX IF NOT EXISTS idx_collaborations_status ON collaborations(status);

CREATE INDEX IF NOT EXISTS idx_shared_content_content_id ON shared_content(content_id);
CREATE INDEX IF NOT EXISTS idx_shared_content_share_url ON shared_content(share_url);
CREATE INDEX IF NOT EXISTS idx_shared_content_expires_at ON shared_content(expires_at);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_content_updated_at BEFORE UPDATE ON content
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_collaborations_updated_at BEFORE UPDATE ON collaborations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_shared_content_updated_at BEFORE UPDATE ON shared_content
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert initial admin user (password: admin123)
INSERT INTO users (email, username, password_hash, first_name, last_name, is_admin, is_verified)
VALUES (
    'admin@open-same.dev',
    'admin',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- admin123
    'System',
    'Administrator',
    TRUE,
    TRUE
) ON CONFLICT (email) DO NOTHING;

-- Insert sample content types
INSERT INTO content (title, description, type, content, creator_id, is_public, is_template)
SELECT 
    'Welcome to Open-Same',
    'A collaborative digital content creation platform',
    'document',
    '# Welcome to Open-Same! 🚀

Open-Same is a comprehensive platform for creating, sharing, and collaborating on digital content. Whether you''re writing documents, coding, or creating diagrams, Open-Same provides the tools you need.

## Features

- **Real-time Collaboration**: Work together with others in real-time
- **Multiple Content Types**: Support for documents, code, diagrams, and more
- **Version Control**: Track changes and maintain history
- **Sharing & Embedding**: Share your work with the world
- **Plugin System**: Extend functionality with custom plugins

## Getting Started

1. Create your first content
2. Invite collaborators
3. Start creating together!

Happy collaborating! 🎉',
    u.id,
    TRUE,
    TRUE
FROM users u WHERE u.username = 'admin'
ON CONFLICT DO NOTHING;

-- Insert sample collaboration
INSERT INTO collaborations (user_id, content_id, role, status)
SELECT 
    u.id,
    c.id,
    'owner',
    'accepted'
FROM users u, content c 
WHERE u.username = 'admin' AND c.title = 'Welcome to Open-Same'
ON CONFLICT DO NOTHING;

-- Create function to get user content with collaborators
CREATE OR REPLACE FUNCTION get_user_content_with_collaborators(user_uuid UUID)
RETURNS TABLE (
    content_id UUID,
    title VARCHAR(255),
    description TEXT,
    type VARCHAR(50),
    is_public BOOLEAN,
    creator_username VARCHAR(100),
    collaboration_role VARCHAR(50),
    collaboration_status VARCHAR(50)
) AS $$
BEGIN
    RETURN QUERY
    SELECT DISTINCT
        c.id as content_id,
        c.title,
        c.description,
        c.type,
        c.is_public,
        u.username as creator_username,
        COALESCE(col.role, 'none') as collaboration_role,
        COALESCE(col.status, 'none') as collaboration_status
    FROM content c
    LEFT JOIN users u ON c.creator_id = u.id
    LEFT JOIN collaborations col ON c.id = col.content_id AND col.user_id = user_uuid
    WHERE c.creator_id = user_uuid OR col.user_id = user_uuid
    ORDER BY c.updated_at DESC;
END;
$$ LANGUAGE plpgsql;

-- Grant permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO opensame;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO opensame;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO opensame;

-- Create a view for public content
CREATE OR REPLACE VIEW public_content AS
SELECT 
    c.id,
    c.title,
    c.description,
    c.type,
    c.created_at,
    u.username as creator_username,
    c.view_count
FROM content c
JOIN users u ON c.creator_id = u.id
WHERE c.is_public = TRUE AND c.deleted_at IS NULL
ORDER BY c.created_at DESC;

-- Grant select on public content view
GRANT SELECT ON public_content TO opensame;