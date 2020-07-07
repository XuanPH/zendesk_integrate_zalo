DROP PROCEDURE IF EXISTS p_zendesk_zalo_mapping_save;

/*
CALL p_erp_applied_object_activate(
	JSON_OBJECT(
		'user_id', 1
        
        ,'saved_by', 1
    )
)
*/
DELIMITER $$
CREATE PROCEDURE p_zendesk_zalo_mapping_save
(
	_uuid varchar(255)
    ,_oaid long
    ,_name varchar(255)
    ,_metadata LONGTEXT
    ,_instance_id varchar(255)
    ,_zendesk_access_token varchar(255)
    ,_subdomain varchar(255)
    ,_locale varchar(255)
    ,_saved_by varchar(255)
)
BEGIN
	IF NOT EXISTS (SELECT id FROM zalo_auth_info where id = _oaid) THEN 
		INSERT INTO `zendesk_integration`.`zalo_auth_info`
		(
                `id`,
				`oaId`,
				`name`,
				`metadata`,
				`instancePushId`,
				`zendeskAccessToken`,
				`subdomain`,
				`locale`,
				`createdAt`,
				`createdBy`,
				`updatedAt`,
				`updatedBy`
		) VALUES
		(
				_uuid,
				_oaid,
				_name,
				_metadata,
				_instance_id,
				_zendesk_access_token,
				_subdomain,
				_locale,
				NOW(),
				_saved_by,
				NOW(),
				_saved_by
		);
    ELSE
		UPDATE `zendesk_integration`.`zalo_auth_info`
		SET
			`name` = _name,
			`metadata` = _metadata,
			`instancePushId` = _instance_id,
			`zendeskAccessToken` = _zendesk_access_token,
			`subdomain` = _subdomain,
			`locale` = _locale,
			`updatedAt` = NOW(),
			`updatedBy` = _saved_by
		WHERE `oaId` = oaid;
	END IF;
END
$$ DELIMITER ;
