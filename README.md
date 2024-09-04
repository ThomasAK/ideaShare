# IdeaShare
A tool for sharing and voting on ideas

### Entities 
- Idea
  - ID
  - Title
  - Description
  - Status
  - CreateBy
  - CreatedAt
  - UpdatedAt
- IdeaComment
  - ID
  - IdeaID
  - Comment
  - CreatedBy
  - CreatedAt
  - UpdatedAt
- IdeaLike
  - ID
  - IdeaID
  - CreatedBy
  - CreatedAt
  - UpdatedAt
    - Unique: IdeaID,CreatedBy
- User
  - ID
  - ExternalID
  - FirstName
  - LastName
- UserRole
  - ID
  - UserID
  - Role
    - Roles: Approver, Admin