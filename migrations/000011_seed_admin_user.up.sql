INSERT INTO admin_user (
    username, 
    email, 
    password, 
    role, 
    is_active
) 
VALUES (
    'admin',                                   
    'admin@gmail.com',                       
    '$2a$10$WwWjFcjiTLmvDuPZ4UI4ku9Wp6N07Kkue3PlZVxWPF/14VDCfNlbG',        
    'super_admin',                            
    1                                     
);