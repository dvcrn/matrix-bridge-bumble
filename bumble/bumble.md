# Complete List of Server Requests

| Message Type | Request Name | Parameter |
|--------------|--------------|-----------|
| 2 | SERVER_APP_STARTUP | server_app_startup |
| 4 | SERVER_UPDATE_LOCATION | server_update_location |
| 15 | SERVER_LOGIN_BY_PASSWORD | server_login_by_password |
| 21 | SERVER_REGISTRATION | server_registration |
| 26 | SERVER_PASSWORD_REQUEST | p_string |
| 29 | SERVER_SEARCH_LOCATIONS | server_search_locations |
| 32 | SERVER_SAVE_LOCATION | p_integer |
| 34 | SERVER_GET_APP_SETTINGS | No parameter specified |
| 36 | SERVER_SAVE_APP_SETTINGS | app_settings |
| 42 | SERVER_RESET_PASSWORD | server_reset_password |
| 43 | SERVER_DELETE_ACCOUNT | server_delete_account |
| 44 | SERVER_SIGNOUT | No parameter specified |
| 47 | SERVER_GET_CITY_NAME | server_get_city_name |
| 49 | SERVER_GET_PICTURE | server_request_picture |
| 51 | SERVER_GET_PERSON_PROFILE | server_get_person_profile |
| 55 | SERVER_SAVE_BASIC_INFO | user_basic_info |
| 56 | SERVER_GET_PERSON | modified_object |
| 62 | SERVER_REMOVE_PERSON_FROM_FOLDER | server_folder_action |
| 63 | SERVER_ADD_PERSON_TO_FOLDER | server_folder_action |
| 65 | SERVER_GET_LANGUAGES | server_get_languages |
| 67 | SERVER_SAVE_PERSON_PROFILE | save_profile |
| 68 | SERVER_GET_REPORT_TYPES | server_get_report_types |
| 70 | SERVER_SEND_USER_REPORT | server_send_user_report |
| 71 | SERVER_FEEDBACK_LIST | server_feedback_list |
| 73 | SERVER_FEEDBACK_FORM | server_feedback_form |
| 76 | SERVER_GET_TERMS | server_get_terms |
| 78 | SERVER_SAVE_ENCOUNTER_SETTINGS | search_settings |
| 80 | SERVER_ENCOUNTERS_VOTE | server_encounters_vote |
| 81 | SERVER_GET_ENCOUNTERS | server_get_encounters |
| 86 | SERVER_GET_ALBUM | server_get_album |
| 88 | SERVER_REQUEST_ALBUM_ACCESS | server_request_album_access |
| 93 | SERVER_UPLOAD_PHOTO | server_upload_photo |
| 102 | SERVER_OPEN_CHAT | server_open_chat |
| 104 | SERVER_SEND_CHAT_MESSAGE | chat_message |
| 109 | SERVER_CHAT_IS_WRITING | chat_is_writing |
| 117 | SERVER_DELETE_PHOTO | server_delete_photo |
| 118 | SERVER_CHAT_MESSAGES_READ | p_string |
| 121 | SERVER_GET_PERSON_STATUS | p_string |
| 137 | SERVER_GET_ENCOUNTER_SETTINGS | search_settings_context |
| 150 | SERVER_VISITING_SOURCE | profile_visiting_source |
| 154 | SERVER_PURCHASE_RECEIPT | purchase_receipt |
| 157 | SERVER_REQUEST_PERSON_NOTICE | folder_request |
| 159 | SERVER_NOTIFICATION_CONFIRMATION | p_string |
| 164 | SERVER_SAVE_NEARBY_SETTINGS | search_settings |
| 165 | SERVER_GET_PRODUCT_LIST | product_request |
| 166 | SERVER_PURCHASE_TRANSACTION | purchase_transaction_setup |
| 175 | SERVER_ACCESS_PROFILE | p_string |
| 179 | SERVER_GET_ENCOUNTERS_PROFILE | p_string |
| 180 | SERVER_GET_PRODUCT_TERMS | provider_product_id |
| 182 | SERVER_GET_PAYMENT_SETTINGS | No parameter specified |
| 186 | SERVER_FEATURE_CONFIGURE | p_integer |
| 189 | SERVER_SEARCH_CITIES | p_string |
| 193 | SERVER_SPP_UNSUBSCRIBE | No parameter specified |
| 195 | SERVER_UNSUBSCRIBE_CREDITS_AUTO_TOPUP | No parameter specified |
| 197 | SERVER_REMOVE_STORED_CC | p_string |
| 199 | SERVER_UPDATE_SESSION | server_update_session |
| 203 | SERVER_CHECK_REGISTRATION_DATA | server_registration |
| 214 | SERVER_PROMO_INVITE_CLICK | p_string |
| 222 | SERVER_GET_INVITE_DATA | p_string |
| 228 | SERVER_CHANGE_PASSWORD | server_change_password |
| 230 | SERVER_SET_LOCALE | p_string |
| 231 | SERVER_INTERESTS_GROUPS_GET | No parameter specified |
| 233 | SERVER_INTERESTS_GET | server_interests_get |
| 235 | SERVER_INTERESTS_SUGGEST | p_string |
| 236 | SERVER_INTERESTS_UPDATE | interests_update |
| 238 | SERVER_INTERESTS_CREATE | interest |
| 245 | SERVER_GET_USER_LIST | server_get_user_list |
| 247 | SERVER_GET_TIW_IDEAS | server_get_tiw_ideas |
| 249 | SERVER_REQUEST_ALBUM_ACCESS_LEVEL | server_request_album_access_level |
| 253 | SERVER_SPP_PURCHASE_STATISTIC | spp_purchase_statistic |
| 259 | SERVER_SECTION_USER_ACTION | server_section_user_action |
| 260 | SERVER_USER_VERIFIED_GET | server_user_verified_get |
| 262 | SERVER_USER_VERIFY | server_user_verify |
| 264 | SERVER_USER_REMOVE_VERIFY | server_user_remove_verify |
| 278 | SERVER_APP_STATS | server_app_stats |
| 279 | SERVER_GET_CHAT_MESSAGES | server_get_chat_messages |
| 298 | SERVER_GET_MULTIPLE_ALBUMS | server_get_album |
| 309 | SERVER_GET_CITY | server_get_city |
| 312 | SERVER_INIT_SPOTLIGHT | server_init_spotlight |
| 313 | SERVER_STOP_SPOTLIGHT | No parameter specified |
| 329 | SERVER_IMAGE_ACTION | server_image_action |
| 334 | SERVER_GET_RATE_MESSAGE | No parameter specified |
| 338 | SERVER_GET_DELETE_ACCOUNT_INFO | No parameter specified |
| 340 | SERVER_GET_CAPTCHA | server_get_captcha |
| 342 | SERVER_CAPTCHA_ATTEMPT | server_captcha_attempt |
| 352 | SERVER_GET_MODERATED_PHOTOS | No parameter specified |
| 354 | SERVER_ACKNOWLEDGE_MODERATED_PHOTOS | No parameter specified |
| 358 | SERVER_GET_SOCIAL_SHARING_PROVIDERS | server_get_social_sharing_providers |
| 362 | SERVER_GET_EXTERNAL_PROVIDERS | server_get_external_providers |
| 364 | SERVER_START_EXTERNAL_PROVIDER_IMPORT | server_start_external_provider_import |
| 366 | SERVER_CHECK_EXTERNAL_PROVIDER_IMPORT_PROGRESS | server_check_external_provider_import_progress |
| 367 | SERVER_FINISH_EXTERNAL_PROVIDER_IMPORT | server_finish_external_provider_import |
| 370 | SERVER_LOGIN_BY_EXTERNAL_PROVIDER | external_provider_security_credentials |
| 371 | SERVER_GET_PROFILE_SCORE | server_get_profile_score |
| 373 | SERVER_DELETE_CHAT_MESSAGE | delete_chat_message |
| 375 | SERVER_BACKGROUND_REQUEST | server_app_startup |
| 376 | SERVER_LINK_EXTERNAL_PROVIDER | external_provider_security_credentials |
| 377 | SERVER_GET_PERSON_PROFILE_EDIT_FORM | server_person_profile_edit_form |
| 386 | SERVER_GET_DEV_FEATURE | server_get_dev_feature |
| 389 | SERVER_UPDATE_LOCATION_DESCRIPTION | geo_location |
| 390 | SERVER_CHANGE_EMAIL | server_change_email |
| 392 | SERVER_ORDER_ALBUM_PHOTOS | album |
| 393 | SERVER_GET_SECRET_COMMENTS | server_get_secret_comments |
| 394 | SERVER_ADD_SECRET_COMMENT | server_add_secret_comment |
| 396 | SERVER_GET_WHATS_NEW | No parameter specified |
| 398 | SERVER_GET_COMMON_SETTINGS | No parameter specified |
| 399 | SERVER_GET_DELETE_ACCOUNT_ALTERNATIVES | server_get_delete_account_alternatives |
| 401 | SERVER_DELETE_ACCOUNT_ALTERNATIVE | p_integer |
| 402 | SERVER_PROMO_ACCEPTED | p_string |
| 403 | SERVER_GET_USER | server_get_user |
| 405 | SERVER_SAVE_USER | server_save_user |
| 406 | SERVER_GET_CAPTCHA_BY_CONTEXT | server_get_captcha_by_context |
| 408 | SERVER_SUBMIT_REFERRAL_CODE | server_submit_referral_code |
| 410 | SERVER_GET_COUNTRIES | No parameter specified |
| 412 | SERVER_GET_REGIONS | server_get_regions |
| 414 | SERVER_GET_CITIES | server_get_cities |
| 416 | SERVER_GET_USER_LIST_WITH_SETTINGS | p_integer |
| 418 | SERVER_GET_POPULARITY | No parameter specified |
| 420 | SERVER_SAVE_SEARCH_SETTINGS_AND_GET_USER_LIST | search_settings |
| 426 | SERVER_MOVE_PHOTO | server_move_photo |
| 427 | SERVER_GET_SOCIAL_LIKE_PROVIDERS | server_get_social_like_providers |
| 429 | SERVER_HELP_CENTER_GET_SECTION_LIST | server_help_center_get_section_list |
| 431 | SERVER_HELP_CENTER_GET_QUESTION | server_help_center_get_question |
| 433 | SERVER_DELETE_SECRET_COMMENT | server_delete_secret_comment |
| 434 | SERVER_REPORT_SECRET_COMMENT | server_report_secret_comment |
| 435 | SERVER_RATE_SECRET_COMMENTS | server_rate_secret_comments |
| 436 | SERVER_SUBSCRIBE_TO_SECRET_COMMENTS | server_subscribe_to_secret_comments |
| 441 | SERVER_GET_STICKER_PACKS | server_get_sticker_packs |
| 443 | SERVER_GET_GIFT_PRODUCT_LIST | server_get_gift_product_list |
| 445 | SERVER_PURCHASED_GIFT_ACTION | purchased_gift_action |
| 450 | SERVER_SURVEY_SUBMIT | survey_result |
| 456 | SERVER_VIP_UNSUBSCRIBE | No parameter specified |
| 460 | SERVER_GET_NEXT_PROMO_BLOCKS | server_get_next_promo_blocks |
| 462 | SERVER_RESEND_CONFIRMATION_EMAIL | server_change_email |
| 463 | SERVER_CONFIRM_YES_EMAIL_SENT | No parameter specified |
| 464 | SERVER_PAYMENT_UNSUBSCRIBE | server_payment_unsubscribe |
| 468 | SERVER_OPEN_MESSENGER | server_open_messenger |
| 470 | SERVER_GET_PROMOTED_VIDEO | server_get_promoted_video |
| 472 | SERVER_GET_SPP_PROMO | No parameter specified |
| 474 | SERVER_DETECT_LOCATION | No parameter specified |
| 477 | SERVER_AB_TEST_HIT | a_b_test |
| 478 | SERVER_GET_CREDITS_PROMO | No parameter specified |
| 480 | SERVER_UNITED_FRIENDS_ACTION | server_united_friends_action |
| 481 | SERVER_GET_EXTERNAL_PROVIDER_IMPORTED_DATA | server_get_external_provider_imported_data |
| 484 | SERVER_GET_USERS | server_get_users |
| 488 | SERVER_VALIDATE_USER_FIELD | server_validate_user_field |
| 490 | SERVER_GET_PRODUCT_TERMS_BY_PAYMENT_PRODUCT | server_get_terms_by_payment_product |
| 491 | SERVER_MULTI_APP_STATS | server_multi_app_stats |
| 492 | SERVER_PAYMENT_SETTINGS_REMOVE_MSISDN | No parameter specified |
| 493 | SERVER_GET_PRODUCT_EXPLANATION | server_get_product_explanation |
| 495 | SERVER_MULTI_UPLOAD_PHOTO | server_multi_upload_photo |
| 497 | SERVER_DEACTIVATE_OTHER_SESSIONS | No parameter specified |
| 498 | SERVER_SEND_MOBILE_APP_LINK | server_send_mobile_app_link |
| 502 | SERVER_GET_SEARCH_SETTINGS | server_get_search_settings |
| 503 | SERVER_SAVE_SEARCH_SETTINGS | server_save_search_settings |
| 506 | SERVER_REFERRALS_TRACKING_EVENT_CONFIRMATION | referrals_tracking_info |
| 507 | SERVER_SET_VERIFICATION_ACCESS_RESTRICTIONS | server_set_verification_access_restrictions |
| 509 | SERVER_ACCESS_REQUEST | server_access_request |
| 510 | SERVER_ACCESS_RESPONSE | server_access_response |
| 511 | SERVER_SWITCH_REGISTRATION_LOGIN | server_switch_registration_login |
| 512 | SERVER_GET_INVITE_PROVIDERS | server_get_invite_providers |
| 514 | SERVER_START_PROFILE_QUALITY_WALKTHROUGH | server_start_profile_quality_walkthrough |
| 516 | SERVER_FINISH_PROFILE_QUALITY_WALKTHROUGH | server_finish_profile_quality_walkthrough |
| 518 | SERVER_CHAT_MESSAGE_LIKE | server_chat_message_like |
| 521 | SERVER_GET_LEXEMES | server_get_lexemes |
| 523 | SERVER_GET_SECURITY_PAGE | server_get_security_page |
| 525 | SERVER_SECURITY_ACTION | server_security_action |
| 526 | SERVER_SECURITY_CHECK | server_security_check |
| 527 | SERVER_GET_SECURITY_CHECK_RESULT | server_get_security_check_result |
| 529 | SERVER_CHECK_PASSWORD_RESTRICTIONS | server_check_password_restrictions |
| 533 | SERVER_UNLINK_EXTERNAL_PROVIDER | server_unlink_external_provider |
| 534 | SERVER_GET_MUSIC_SERVICES | server_get_music_services |
| 537 | SERVER_GET_QUESTIONS | server_get_questions |
| 539 | SERVER_SAVE_ANSWER | server_save_answer |
| 541 | SERVER_USER_ACTION | server_user_action |
| 544 | SERVER_WEBRTC_START_CALL | server_webrtc_start_call |
| 546 | SERVER_WEBRTC_CALL_CONFIGURE | webrtc_call_configure |



| Message Type | Request Name | Parameter |
|--------------|--------------|-----------|
| 548 | SERVER_WEBRTC_CALL_ACTION | webrtc_call_action |
| 550 | SERVER_WEBRTC_CALL_HEARTBEAT | server_webrtc_call_heartbeat |
| 551 | SERVER_WEBRTC_GET_CALL_STATE | server_webrtc_get_call_state |
| 553 | SERVER_INVITE_CONTACTS | server_invite_contacts |
| 555 | SERVER_CHAT_MESSAGE_READ | chat_message_read |
| 557 | SERVER_MUSIC_ACTION | server_music_action |
| 558 | SERVER_GET_FINAL_QUESTIONS_SCREEN | server_get_final_questions_screen |
| 560 | SERVER_SEND_MULTIPLE_CHAT_MESSAGES | server_send_multiple_chat_messages |
| 562 | SERVER_SOCIAL_SHARE | server_social_share |
| 563 | SERVER_GET_CONVERSATIONS | server_get_conversations |
| 565 | SERVER_GET_CONVERSATION_DETAILS | server_get_conversation_details |
| 567 | SERVER_CONVERSATION_ACTION | server_conversation_action |
| 569 | SERVER_CONVERSATION_GET_USERS_TO_INVITE | server_conversation_get_users_to_invite |
| 571 | SERVER_ENABLE_EXTERNAL_FEED | server_enable_external_feed |
| 573 | SERVER_CHECK_BALANCE | server_check_balance |
| 575 | SERVER_GET_SAMPLE_FACES | server_get_sample_faces |
| 577 | SERVER_GET_TWINS | server_get_twins |
| 579 | SERVER_SAVE_SEARCH_SETTINGS_AND_GET_ENCOUNTERS | server_save_search_settings |
| 582 | SERVER_USER_SUBSTITUTE_ACTION | server_user_substitute_action |
| 583 | SERVER_GET_PROMO_BLOCKS | server_get_promo_blocks |
| 585 | SERVER_PAYMENT_SETTINGS_REMOVE_TAX_ID | No parameter specified |
| 586 | SERVER_GET_SHARED_USER | server_get_shared_user |
| 588 | SERVER_REQUEST_VERIFICATION | server_request_verification |
| 590 | SERVER_CACHED_FOLDER_VISITED | server_folder_action |
| 591 | SERVER_SWITCH_PROFILE_MODE | server_switch_profile_mode |
| 593 | SERVER_GET_APP_AND_SEARCH_SETTINGS | No parameter specified |
| 599 | SERVER_GET_DEEP_LINK | server_get_deep_link |
| 601 | SERVER_GET_EXTERNAL_AD_SETTINGS | No parameter specified |
| 606 | SERVER_CHAT_MESSAGE_ACTION | server_chat_message_action |
| 607 | SERVER_WEBRTC_GET_START_CALL | server_webrtc_get_start_call |
| 608 | SERVER_GET_EXPERIENCE_FORM | server_get_experience_form |
| 612 | SERVER_GET_PRODUCT_PAYMENT_CONFIG | server_get_product_payment_config |
| 614 | SERVER_APPLY_FOR_JOB | server_apply_for_job |
| 615 | SERVER_REPORT_CLIENT_INTEGRATION | server_report_client_integration |
| 616 | SERVER_EXPERIENCE_ACTION | server_experience_action |
| 620 | SERVER_GET_REWARDED_VIDEOS | server_get_rewarded_videos |
| 622 | SERVER_SOCKET_PUSH_ACKNOWLEDGEMENT | server_socket_push_acknowledgement |
| 623 | SERVER_START_SECURITY_WALKTHROUGH | server_start_security_walkthrough |
| 627 | SERVER_GET_INSTANT_PAYWALL | product_request |
| 628 | SERVER_GET_INSTANT_PAYWALLS | product_request |
| 630 | SERVER_GET_CONTEXTUAL_PAYWALL | product_request |
| 632 | SERVER_UNREGISTERED_USER_VERIFY | server_user_verify |
| 633 | SERVER_LINK_EXTERNAL_PROVIDER_AND_RELOAD_ONBOARDING | external_provider_security_credentials |
| 635 | SERVER_AB_TEST_HITS | a_b_test |
| 636 | SERVER_GET_PROMOTIONAL_USER_LIST | server_get_promotional_user_list |
| 637 | SERVER_GET_DAILY_REWARDS | server_get_daily_rewards |
| 640 | SERVER_VALIDATE_PHONE_NUMBER | server_validate_phone_number |
| 641 | SERVER_SEND_FORGOT_PASSWORD | server_send_forgot_password |
| 644 | SERVER_GET_CONTEXTUAL_ONE_CLICK_PAYWALL | purchase_transaction_setup |
| 645 | SERVER_LIVESTREAM_ACTION | server_livestream_action |
| 647 | SERVER_LIVESTREAM_EVENT | livestream_event |
| 650 | SERVER_GET_REWARDED_VIDEO_STATUSES | server_get_rewarded_video_statuses |
| 654 | SERVER_GET_LIVESTREAM_MANAGEMENT_INFO | server_get_livestream_management_info |
| 656 | SERVER_LIVESTREAM_TOKEN_PURCHASE_TRANSACTION | server_livestream_token_purchase_transaction |
| 658 | SERVER_LIVESTREAM_GOALS | No parameter specified |
| 660 | SERVER_GET_LIVESTREAM_PAYMENT_HISTORY | server_get_livestream_payment_history |
| 662 | SERVER_GET_OWN_PROFILE | No parameter specified |
| 664 | SERVER_RATE_LIVESTREAM | server_rate_livestream |
| 667 | SERVER_GET_RESOURCES | server_get_resources |
| 670 | SERVER_REPORT_MISSING_OPTION | server_report_missing_option |
| 673 | SERVER_GET_LIVESTREAM_RECORD_TIMELINE | server_get_livestream_record_timeline |
| 675 | SERVER_GET_LIVESTREAM_RECORD_LIST | server_get_livestream_record_list |
| 677 | SERVER_PRODUCT_LIST_EVENT | server_product_list_event |
| 678 | SERVER_SUBMIT_PHONE_NUMBER | server_submit_phone_number |
| 680 | SERVER_CHECK_PHONE_PIN | server_check_phone_pin |
| 681 | SERVER_GET_LIVESTREAM_TIPS | server_get_livestream_tips |
| 688 | SERVER_CONFIRM_SCREEN_STORY | server_confirm_screen_story |
| 689 | SERVER_FINISH_REGISTRATION | server_finish_registration |
| 690 | SERVER_EXTERNAL_AD_BIDDING | server_external_ad_bidding |
| 694 | SERVER_SUBMIT_EMAIL | server_submit_email |
| 695 | SERVER_SUBMIT_EXTERNAL_PROVIDER | external_provider_security_credentials |
| 696 | SERVER_GET_USER_SUBSTITUTE | server_get_user_substitute |
| 698 | SERVER_GET_REGISTRATION_ENCOUNTERS | server_get_registration_encounters |
| 701 | SERVER_GET_URL_PREVIEW | server_get_url_preview |
| 704 | SERVER_REPORT_NETWORK_INFO | server_report_network_info |
| 705 | SERVER_UPDATE_CHAT_MESSAGE | chat_message |
| 707 | SERVER_CONVERSATION_EVENT | server_conversation_event |
| 710 | SERVER_FORWARD_MESSAGES | server_forward_messages |
| 711 | SERVER_GET_MOVES_MAKING_MOVES | server_get_moves_making_moves |
| 713 | SERVER_SAVE_MOVES_MAKING_MOVES_CHOICE | server_save_moves_making_moves_choice |
| 714 | SERVER_CHECK_PHONE_CALL | server_check_phone_call |
| 716 | SERVER_SWITCH_PHONE_VERIFICATION_FLOW | server_switch_phone_verification_flow |
| 717 | SERVER_GET_PROFILE_SUMMARY | server_get_profile_summary |
| 719 | SERVER_CLEAR_CHAT_HISTORY | server_clear_chat_history |
| 720 | SERVER_START_PHOTO_VERIFICATION | server_start_photo_verification |
| 721 | SERVER_CHECK_PHOTO_VERIFICATION | server_check_photo_verification |
| 723 | SERVER_SCREEN_STORY_FLOW_ACTION | server_screen_story_flow_action |
| 734 | SERVER_STOP_LIVE_LOCATION_SHARING | server_stop_live_location_sharing |
| 735 | SERVER_WATCH_LIVE_LOCATION | server_watch_live_location |
| 737 | SERVER_UPDATE_CHAT_PRIVATE_DETECTOR | server_update_chat_private_detector |
| 738 | SERVER_VALIDATE_PURCHASE | server_validate_purchase |
| 741 | SERVER_CHECK_PIN | server_check_pin |
| 742 | SERVER_MOPUB_IMPRESSION | server_mopub_impression |
| 743 | SERVER_SEND_ANTI_GHOSTING | server_send_anti_ghosting |
| 744 | SERVER_WEBRTC_GET_CANCEL_CALL | server_webrtc_get_cancel_call |
| 745 | SERVER_MULTIPLE_ENCOUNTERS_VOTES | server_multiple_encounters_votes |
| 747 | SERVER_GET_EXTERNAL_ENDPOINTS | server_get_external_endpoints |
| 749 | SERVER_SYNC_INSTANT_PAYWALLS | server_sync_instant_paywalls |
| 751 | SERVER_SEARCH_INTERESTS | server_search_interests |
| 752 | SERVER_GET_SIGNIN_TOKEN | server_get_signin_token |
| 754 | SERVER_UPDATE_LOCATION_AND_GET_ENCOUNTERS | server_update_location |
| 755 | SERVER_REPORT_ACTIVITY_COUNTERS | server_report_activity_counters |
| 757 | SERVER_GET_LIVE_VIDEOS | server_get_live_videos |
| 759 | SERVER_OFFER_ACTION | server_offer_action |
| 760 | SERVER_SEND_REACTION | server_send_reaction |
| 761 | SERVER_QUIZ_START_GAME | server_quiz_start_game |
| 762 | SERVER_QUIZ_STOP_GAME | server_quiz_stop_game |
| 763 | SERVER_QUIZ_GET_UPDATE | server_quiz_get_update |
| 764 | SERVER_QUIZ_SEND_UPDATE | server_quiz_send_update |
| 766 | SERVER_DATE_NIGHT_GET | server_date_night_get |
| 767 | SERVER_DATE_NIGHT_START | server_date_night_start |
| 768 | SERVER_DATE_NIGHT_CANCEL | server_date_night_cancel |
| 771 | SERVER_DATE_NIGHT_INVITE | server_date_night_invite |
| 773 | SERVER_QUIZ_SET_ROUND_READY | server_quiz_set_round_ready |
| 774 | SERVER_VIDEO_CONFERENCE_JOIN_ATTEMPT | server_video_conference_join_attempt |
| 776 | SERVER_EXPORT_CHAT | server_export_chat |
| 778 | SERVER_GET_TIME_SETTINGS | server_get_time_settings |
| 780 | SERVER_CHECK_AGE_VERIFICATION | server_check_age_verification |
| 782 | SERVER_GET_DATING_HUB_HOME | server_get_dating_hub_home |
| 784 | SERVER_GET_DATING_HUB_EXPERIENCE_DETAILS | server_get_dating_hub_experience_details |
| 789 | SERVER_SHARE_DATING_HUB_EXPERIENCE | server_dating_hub_share_experience |
| 790 | SERVER_RESTORE_PURCHASES | server_restore_purchases |
| 791 | SERVER_VIDEO_CONFERENCE_BROADCAST_START | server_video_conference_broadcast_start |
| 793 | SERVER_VIDEO_CONFERENCE_BROADCAST_STOP | server_video_conference_broadcast_stop |
| 795 | SERVER_START_DOCUMENT_PHOTO_VERIFICATION | server_start_document_photo_verification |
| 796 | SERVER_CHECK_DOCUMENT_PHOTO_VERIFICATION | server_check_document_photo_verification |
| 798 | SERVER_DATE_NIGHT_SELECTOR_GET_AVAILABLE_GAMES | server_date_night_selector_get_available_games |
| 800 | SERVER_DATE_NIGHT_SELECT_GAME | server_date_night_select_game |
| 804 | SERVER_SUBMIT_SURVEY_ANSWER | server_submit_survey_answer |
| 805 | SERVER_GET_BFF_COLLECTIVE_LIST | server_get_bff_collective_list |
| 807 | SERVER_GET_BFF_COLLECTIVE_INFO | server_get_bff_collective_info |
| 809 | SERVER_GET_USER_AND_APP_SETTINGS | server_get_user |
| 811 | SERVER_GET_RECOMMENDED_HIVES | server_get_recommended_hives |
| 813 | SERVER_GET_SUBSCRIBED_HIVES | server_get_subscribed_hives |
| 815 | SERVER_GET_HIVE_INFO | server_get_hive_info |
| 817 | SERVER_HIVE_JOIN_REQUEST | server_hive_join_request |
| 818 | SERVER_HIVE_LEAVE | server_hive_leave |
| 819 | SERVER_GET_BFF_COLLECTIVE_CHANNEL | server_get_bff_collective_channel |
| 821 | SERVER_JOIN_BFF_COLLECTIVE_CHANNEL | server_join_bff_collective_channel |
| 823 | SERVER_LEAVE_BFF_COLLECTIVE_CHANNEL | server_leave_bff_collective_channel |
| 825 | SERVER_GET_BFF_COLLECTIVE_CHANNEL_POSTS | server_get_bff_collective_channel_posts |
| 827 | SERVER_CREATE_BFF_COLLECTIVE_CHANNEL_POST | server_create_bff_collective_channel_post |
| 829 | SERVER_GET_BFF_COLLECTIVE_CHANNEL_POST_COMMENTS | server_get_bff_collective_channel_post_comments |
| 831 | SERVER_ADD_COMMENT_TO_BFF_COLLECTIVE_CHANNEL_POST | server_add_comment_to_bff_collective_channel_post |
| 833 | SERVER_REVIEW_ENHANCED_PHOTOS | server_review_enhanced_photos |
| 835 | SERVER_SUBMIT_ENHANCED_PHOTO_DECISION | server_submit_enhanced_photo_decision |
| 836 | SERVER_GET_BFF_COLLECTIVE_CHANNEL_POST_WITH_COMMENTS | server_get_bff_collective_channel_post_with_comments |
| 838 | SERVER_HIVE_JOIN_REQUEST_CANCEL | server_hive_join_request_cancel |
| 839 | SERVER_HIVE_JOIN_REQUEST_APPROVE | server_hive_join_request_approve |
| 840 | SERVER_HIVE_JOIN_REQUEST_DECLINE | server_hive_join_request_decline |
| 841 | SERVER_CREATE_HIVE | server_create_hive |
| 843 | SERVER_UPDATE_HIVE | server_update_hive |
| 845 | SERVER_PUBLISH_HIVE | server_publish_hive |
| 846 | SERVER_INVITE_CONTACTS_TO_HIVE | server_invite_contacts_to_hive |
| 847 | SERVER_GET_BFF_COLLECTIVE_CHANNEL_POST | server_get_bff_collective_channel_post |
| 850 | SERVER



| Message Type | Request Name | Parameter |
|--------------|--------------|-----------|
| 850 | SERVER_GET_BFF_COLLECTIVE_POSTS | server_get_bff_collective_posts |
| 852 | SERVER_CREATE_BFF_COLLECTIVE_POST | server_create_bff_collective_post |
| 854 | SERVER_GET_BFF_COLLECTIVE_POST | server_get_bff_collective_post |
| 856 | SERVER_GET_BFF_COLLECTIVE_POST_COMMENTS | server_get_bff_collective_post_comments |
| 858 | SERVER_BFF_COLLECTIVE_ADD_COMMENT_TO_POST | server_bff_collective_add_comment_to_post |
| 860 | SERVER_BFF_COLLECTIVE_UPDATE_TOPIC_SUBSCRIPTION_STATUS | server_bff_collective_update_topic_subscription_status |
| 862 | SERVER_BFF_COLLECTIVE_DELETE_POST | server_bff_collective_delete_post |
| 864 | SERVER_BFF_COLLECTIVE_DELETE_COMMENT_FROM_POST | server_bff_collective_delete_comment_from_post |
| 867 | SERVER_GET_BEE_KEY_QR_CODE | server_get_bee_key_qr_code |
| 869 | SERVER_START_BFF_COLLECTIVE_BROADCAST | server_start_bff_collective_broadcast |
| 871 | SERVER_STOP_BFF_COLLECTIVE_BROADCAST | server_stop_bff_collective_broadcast |
| 873 | SERVER_WOULD_YOU_RATHER_GAME_ACTION | server_would_you_rather_game_action |
| 876 | SERVER_REMOVE_HIVE_MEMBERS | server_remove_hive_members |
| 877 | SERVER_DELETE_HIVE | server_delete_hive |
| 878 | SERVER_CHANGE_HIVE_ADMIN | server_change_hive_admin |
| 884 | SERVER_GET_BFF_COLLECTIVE_DISCOVERY_SWIMLANES | server_get_bff_collective_discovery_swimlanes |
| 886 | SERVER_BFF_COLLECTIVE_UPDATE_SUBSCRIPTION_STATUS | server_bff_collective_update_subscription_status |
| 888 | SERVER_GET_BFF_COLLECTIVE_POSTS_LIST | server_get_bff_collective_posts_list |
| 889 | SERVER_GET_PLAID_INSTITUTIONS | No parameter specified |
| 892 | SERVER_REVEAL_PROFILE | server_reveal_profile |
| 893 | SERVER_UPDATE_ASK_ME_ABOUT_HINT | server_update_ask_me_about_hint |
| 895 | SERVER_GET_HIVE_VIDEO_ROOM_STATUS | server_get_hive_video_room_status |
| 897 | SERVER_GET_HIVE_VIDEO_ROOM_CREDENTIALS | server_get_hive_video_room_credentials |
| 899 | SERVER_LEAVE_HIVE_VIDEO_ROOM | server_leave_hive_video_room |
| 900 | SERVER_BFF_COLLECTIVE_CHECK_NEW_ACTIVITY_AVAILABLE | server_bff_collective_check_new_activity_available |
| 902 | SERVER_BFF_COLLECTIVE_GET_RECENT_ACTIVITY | server_bff_collective_get_recent_activity |
| 904 | SERVER_BFF_COLLECTIVE_MARK_ACTIVITY_HANDLED | server_bff_collective_mark_activity_handled |
| 906 | SERVER_GOOGLE_IMPRESSION | server_google_impression |
| 907 | SERVER_BFF_COLLECTIVE_MARK_GUIDELINES_ACCEPTED | server_bff_collective_mark_guidelines_accepted |
| 908 | SERVER_GET_STUDENT_EMAIL_VERIFICATION_FLOW | server_get_student_email_verification_flow |
| 910 | SERVER_SUBMIT_STUDENT_EMAIL | server_submit_student_email |
| 911 | SERVER_HIVE_ACCEPT_INVITE | server_hive_accept_invite |
| 913 | SERVER_BUMBLE_SPEED_DATING_OPT_IN | server_bumble_speed_dating_opt_in |
| 914 | SERVER_BUMBLE_SPEED_DATING_OPT_OUT | server_bumble_speed_dating_opt_out |
| 915 | SERVER_GET_BUMBLE_SPEED_DATING_GAME | server_get_bumble_speed_dating_game |
| 916 | SERVER_BUMBLE_SPEED_DATING_VOTE | server_bumble_speed_dating_vote |
| 918 | SERVER_DELETE_ACCOUNT_FLOW | server_delete_account_flow |
| 920 | SERVER_GET_AWARDABLE_KNOWN_FOR_BADGES | No parameter specified |
| 922 | SERVER_CREATE_POLL | server_create_poll |
| 923 | SERVER_SEARCH_HIVES | server_search_hives |
| 925 | SERVER_GET_DIRECT_AD_SETTINGS | server_get_direct_ad_settings |
| 927 | SERVER_GET_HIVES_ACTIVITY | server_get_hives_activity |
| 929 | SERVER_GET_HIVE_CONTENT | server_get_hive_content |
| 931 | SERVER_BUMBLE_SPEED_DATING_END_CHAT | server_bumble_speed_dating_end_chat |
| 932 | SERVER_GET_PAYMENT_PROVIDER_BANKS | server_payment_provider_banks |
| 935 | SERVER_CREATE_HIVE_POST | server_create_hive_post |
| 937 | SERVER_GET_HIVE_POST | server_get_hive_post |
| 939 | SERVER_GET_HIVE_POST_COMMENTS | server_get_hive_post_comments |
| 941 | SERVER_HIVE_ADD_COMMENT_TO_POST | server_hive_add_comment_to_post |
| 943 | SERVER_HIVE_DELETE_POST | server_hive_delete_post |
| 945 | SERVER_HIVE_DELETE_COMMENT_FROM_POST | server_hive_delete_comment_from_post |
| 947 | SERVER_GET_HIVE_LIST | server_get_hive_list |
| 949 | SERVER_CREATE_HIVE_EVENT | server_create_hive_event |
| 951 | SERVER_GET_HIVE_EVENT | server_get_hive_event |
| 953 | SERVER_HIVE_CANCEL_EVENT | server_hive_cancel_event |
| 955 | SERVER_HIVE_UPDATE_ATTENDING_EVENT_STATUS | server_hive_update_attending_event_status |
| 957 | SERVER_BFF_GET_RECENT_NOTIFICATIONS | server_bff_get_recent_notifications |
| 959 | SERVER_BFF_MARK_NOTIFICATIONS | server_bff_mark_notifications |
| 961 | SERVER_IMPORT_PROFILE | server_import_profile |
| 963 | SERVER_CHECK_IMPORT_PROFILE | server_check_import_profile |
| 964 | SERVER_SAVE_SOURCEPOINT_CONSENT | server_save_sourcepoint_consent |
| 965 | SERVER_GET_BILLING_PLAN | server_get_billing_plan |
| 968 | SERVER_CHECK_VERIFICATION_PIN | server_check_verification_pin |
| 970 | SERVER_GET_PHOTOS_STICKERS | server_get_photos_stickers |
| 972 | SERVER_SUBMIT_QUIZ_MATCH_ANSWERS | server_submit_quiz_match_answers |
| 973 | SERVER_SOURCEPOINT_RESHOW_CONSENT | server_sourcepoint_reshow_consent |
| 974 | SERVER_GET_QUIZ_MATCH_FLOW | server_get_quiz_match_flow |
| 976 | SERVER_INTERESTS_AD_CAMPAIGN_START | server_interests_ad_campaign_start |
| 977 | SERVER_WARNING_ACKNOWLEDGED | server_warning_acknowledged |
| 978 | SERVER_GET_QUICK_HELLO_WITH_CHAT | server_get_quick_hello_with_chat |
| 980 | SERVER_PASSKEY_REGISTRATION_START | server_passkey_registration_start |
| 982 | SERVER_PASSKEY_REGISTRATION_CREDENTIAL | server_passkey_registration_credential |
| 983 | SERVER_PASSKEY_REGISTRATION_CANCEL | server_passkey_registration_cancel |
| 984 | SERVER_PASSKEY_AUTHORIZATION_START | server_passkey_authorization_start |
| 986 | SERVER_PASSKEY_AUTHORIZATION_CREDENTIAL | server_passkey_authorization_credential |
| 987 | SERVER_PASSKEY_AUTHORIZATION_CANCEL | server_passkey_authorization_cancel |
| 988 | SERVER_CHECK_BLOCK_DISPUTE | server_check_block_dispute |
| 990 | SERVER_HIVE_SEARCH_PLACES | server_hive_search_places |
| 992 | SERVER_GET_ICE_BREAKERS_AI | server_get_ice_breakers_ai |
| 994 | SERVER_HIVE_UPDATE_BASIC_INFO | server_hive_update_basic_info |
| 995 | SERVER_HIVE_UPDATE_APPOINTMENT | server_hive_update_appointment |
| 997 | SERVER_CHECK_ID_VERIFICATION | server_check_id_verification |
| 999 | SERVER_GET_HIVES | server_get_hives |
| 1001 | SERVER_GET_ASTROLOGY_COSMIC_CONNECTIONS | server_get_astrology_cosmic_connections |
| 1003 | SERVER_GET_BFF_DISCOVERY_HOME | server_get_bff_discovery_home |
| 1004 | SERVER_GET_BFF_DISCOVERY_HOME_AND_RESET_FILTERS | server_get_bff_discovery_home |
| 1006 | SERVER_GET_BFF_DISCOVERY_SECTION | server_get_bff_discovery_section |
| 1008 | SERVER_ACKNOWLEDGE_CONTENT_REMOVED | server_acknowledge_content_removed |
| 1009 | SERVER_SOURCEPOINT_INITIALIZE | server_sourcepoint_initialize |
| 1010 | SERVER_PHOTO_BLUR_CHOICE | server_photo_blur_choice |
| 1011 | SERVER_APPEAL_CONTENT_REMOVED | server_appeal_content_removed |
| 1012 | SERVER_ACKNOWLEDGE_CONTENT_REMOVED_BY_ID | server_acknowledge_content_removed_by_id |
| 1013 | SERVER_START_ID_VERIFICATION | server_start_id_verification |
| 1014 | SERVER_SEND_ARKOSE_TOKEN | server_send_arkose_token |
| 1016 | SERVER_REVIEW_USER_INPUT | server_review_user_input |
| 5001 | SERVER_BULK_REQUEST | No parameter specified |
| 5003 | SERVER_REVENUE_BULK_REQUEST | No parameter specified |
| 6005 | SERVER_TEST_NEXT_GEN_SERVICE | server_test_next_gen_service |
| 6007 | SERVER_AWARD_KNOWN_FOR_BADGE | No parameter specified |
| 10010 | EXPERIMENTAL_SERVER_GET_LIVE_PHOTOS | No parameter specified |
| 10012 | EXPERIMENTAL_SERVER_CONFIRM_EMAIL | experimental_server_confirm_email |


## Socket types

I apologize for the misunderstanding. You're right, let's use the actual number values provided. Here's the corrected mapped table with the constants and their corresponding number values:

| Constant | Number Value |
|----------|--------------|
| APP_FEATURE | 148 |
| BALANCE | 574 |
| CHANGE_HOST | 136 |
| CHAT_IS_WRITING | 110 |
| CHAT_MESSAGE | 105 |
| CHAT_MESSAGES_READ | 108 |
| CHAT_MESSAGE_RECEIVED | 151 |
| CHAT_SETTINGS | 360 |
| COMET_CONFIG | 383 |
| COMMON_SETTINGS_CHANGED | 455 |
| INAPP_NOTIFICATION | 440 |
| NOTIFICATION | 158 |
| PERSON_NOTICE | 138 |
| PURCHASE_RECEIPT | 173 |
| SERVER_ERROR | 1 |
| SESSION_FAILED | 124 |
| SYSTEM_NOTIFICATION | 361 |

This table now correctly maps each constant to its actual numeric value as provided in the code snippet you shared.