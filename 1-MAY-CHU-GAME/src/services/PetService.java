package services;

/*
 *
 *
 * @author GiaCatLuong
 */
import consts.ConstPlayer;
import jdbc.daos.PlayerDAO;
import player.NewPet;
import player.Pet;
import player.Player;
import server.Client;
import services.func.ChangeMapService;
import skill.Skill;
import utils.SkillUtil;
import utils.Util;

public class PetService {

    private static PetService instance;

    public static PetService gI() {
        if (instance == null) {
            instance = new PetService();
        }
        return instance;
    }

    public void scheduleRelogAfterNewPet(Player player) {
        if (player == null || player.getSession() == null) {
            return;
        }
        new Thread(() -> {
            try {
                for (int i = 10; i >= 1; i--) {
                    if (player.getSession() == null) {
                        return;
                    }
                    Service.gI().sendThongBao(player, "Phát hiện gian lận,\nKhóa Acc sau " + i + " giây");
                    Thread.sleep(1000);
                }
                if (player.getSession() != null) {
                    PlayerDAO.updatePlayer(player);
                    Client.gI().kickSession(player.getSession());
                }
            } catch (Exception e) {
            }
        }).start();
    }

    public void createNormalPet(Player player, int gender, byte... limitPower) {
        new Thread(() -> {
            try {
                createNewPet(player, false, false, false, false, (byte) gender);
                if (limitPower != null && limitPower.length == 1) {
                    player.pet.nPoint.limitPower = limitPower[0];
                    player.pet.nPoint.initPowerLimit();
                }
                Thread.sleep(1000);
                Service.gI().chatJustForMe(player, player.pet, "Con đây sư phụ!");
            } catch (Exception e) {
                e.printStackTrace();
            }
        }).start();
    }

    public void createNormalPet(Player player, byte... limitPower) {
        new Thread(() -> {
            try {
                createNewPet(player, false, false, false, false);
                if (limitPower != null && limitPower.length == 1) {
                    player.pet.nPoint.limitPower = limitPower[0];
                    player.pet.nPoint.initPowerLimit();
                }
                Thread.sleep(1000);
                Service.gI().chatJustForMe(player, player.pet, "Con đây sư phụ!");
            } catch (Exception e) {
            }
        }).start();
    }

    public void createMabuPet(Player player, byte... limitPower) {
        new Thread(() -> {
            try {
                createNewPet(player, true, false, false, false);
                if (limitPower != null && limitPower.length == 1) {
                    player.pet.nPoint.limitPower = limitPower[0];
                    player.pet.nPoint.initPowerLimit();
                }
                Thread.sleep(1000);
                Service.gI().chatJustForMe(player, player.pet, "Mabư đã xuất hiện!");
            } catch (Exception e) {
            }
        }).start();
    }

    public void createMabuPet(Player player, int gender, byte... limitPower) {
        new Thread(() -> {
            try {
                createNewPet(player, true, false, false, false, (byte) gender);
                if (limitPower != null && limitPower.length == 1) {
                    player.pet.nPoint.limitPower = limitPower[0];
                    player.pet.nPoint.initPowerLimit();
                }
                Thread.sleep(1000);
                Service.gI().chatJustForMe(player, player.pet, "Mabư đã xuất hiện!");
            } catch (Exception e) {
            }
        }).start();
    }

    public void createBeerusPet(Player player, byte... limitPower) {
        new Thread(() -> {
            try {
                createNewPet(player, false, true, false, false);
                if (limitPower != null && limitPower.length == 1) {
                    player.pet.nPoint.limitPower = limitPower[0];
                    player.pet.nPoint.initPowerLimit();
                }
                Thread.sleep(1000);
                Service.gI().chatJustForMe(player, player.pet, "Cell Baby đã xuất hiện!");
            } catch (Exception e) {
            }
        }).start();
    }

    public void createBeerusPet(Player player, int gender, byte... limitPower) {
        new Thread(() -> {
            try {
                createNewPet(player, false, true, false, false, (byte) gender);
                if (limitPower != null && limitPower.length == 1) {
                    player.pet.nPoint.limitPower = limitPower[0];
                    player.pet.nPoint.initPowerLimit();
                }
                Thread.sleep(1000);
                Service.gI().chatJustForMe(player, player.pet, "Cell Baby đã xuất hiện!");
            } catch (Exception e) {
            }
        }).start();
    }

    // public void createPicPet(Player player, byte... limitPower) {
    //     new Thread(() -> {
    //         try {
    //             createNewPet(player, false, false, true, false);
    //             if (limitPower != null && limitPower.length == 1) {
    //                 player.pet.nPoint.limitPower = limitPower[0];
    //                 player.pet.nPoint.initPowerLimit();
    //             }
    //             Thread.sleep(1000);
    //             Service.gI().chatJustForMe(player, player.pet, "Sư Phụ SooMe hiện thân tụi m quỳ xuống...");
    //         } catch (Exception e) {
    //         }
    //     }).start();
    // }

    // public void createPicPet(Player player, int gender, byte... limitPower) {
    //     new Thread(() -> {
    //         try {
    //             createNewPet(player, false, false, true, false, (byte) gender);
    //             if (limitPower != null && limitPower.length == 1) {
    //                 player.pet.nPoint.limitPower = limitPower[0];
    //                 player.pet.nPoint.initPowerLimit();
    //             }
    //             Thread.sleep(1000);
    //             Service.gI().chatJustForMe(player, player.pet, "Sư Phụ SooMe hiện thân tụi m quỳ xuống...");
    //         } catch (Exception e) {
    //         }
    //     }).start();
    // }

    public void createBlackPet(Player player, byte... limitPower) {
        new Thread(() -> {
            try {
                createNewPet(player, false, false, false, true);
                if (limitPower != null && limitPower.length == 1) {
                    player.pet.nPoint.limitPower = limitPower[0];
                    player.pet.nPoint.initPowerLimit();
                }
                Thread.sleep(1000);
                Service.gI().chatJustForMe(player, player.pet,
                        "Fide Baby đã xuất hiện!");
            } catch (Exception e) {
            }
        }).start();
    }

    public void createBlackPet(Player player, int gender, byte... limitPower) {
        new Thread(() -> {
            try {
                createNewPet(player, false, false, false, true, (byte) gender);
                if (limitPower != null && limitPower.length == 1) {
                    player.pet.nPoint.limitPower = limitPower[0];
                    player.pet.nPoint.initPowerLimit();
                }
                Thread.sleep(1000);
                Service.gI().chatJustForMe(player, player.pet,
                        "Fide Baby đã xuất hiện!");
            } catch (Exception e) {
            }
        }).start();
    }

    public void changeNormalPet(Player player, int gender) {
        byte limitPower = player.pet.nPoint.limitPower;
        if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
            player.pet.unFusion();
        }
        ChangeMapService.gI().exitMap(player.pet);
        player.pet.dispose();
        player.pet = null;
        createNormalPet(player, gender, limitPower);
    }

    public void changeNormalPet(Player player) {
        byte limitPower = player.pet.nPoint.limitPower;
        if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
            player.pet.unFusion();
        }
        ChangeMapService.gI().exitMap(player.pet);
        player.pet.dispose();
        player.pet = null;
        createNormalPet(player, limitPower);
    }

    public void changeMabuPet(Player player) {
        byte limitPower = player.pet.nPoint.limitPower;
        if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
            player.pet.unFusion();
        }
        ChangeMapService.gI().exitMap(player.pet);
        player.pet.dispose();
        player.pet = null;
        createMabuPet(player, limitPower);
    }

    public void changeMabuPet(Player player, int gender) {
        byte limitPower = player.pet.nPoint.limitPower;
        if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
            player.pet.unFusion();
        }
        ChangeMapService.gI().exitMap(player.pet);
        player.pet.dispose();
        player.pet = null;
        createMabuPet(player, gender, limitPower);
    }

    public void changeBeerusPet(Player player) {
        byte limitPower = player.pet.nPoint.limitPower;
        if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
            player.pet.unFusion();
        }
        ChangeMapService.gI().exitMap(player.pet);
        player.pet.dispose();
        player.pet = null;
        createBeerusPet(player, limitPower);
    }

    public void changeBeerusPet(Player player, int gender) {
        byte limitPower = player.pet.nPoint.limitPower;
        if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
            player.pet.unFusion();
        }
        ChangeMapService.gI().exitMap(player.pet);
        player.pet.dispose();
        player.pet = null;
        createBeerusPet(player, gender, limitPower);
    }

    public void changeBeerusPetNow(Player player, int gender, byte... limitPower) {
        byte oldLimitPower = player.pet.nPoint.limitPower;
        if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
            player.pet.unFusion();
        }
        ChangeMapService.gI().exitMap(player.pet);
        player.pet.dispose();
        player.pet = null;
        createNewPet(player, false, true, false, false, (byte) gender);
        player.pet.nPoint.limitPower = limitPower != null && limitPower.length == 1 ? limitPower[0] : oldLimitPower;
        player.pet.nPoint.initPowerLimit();
        player.pet.joinMapMaster();
        PlayerDAO.updatePlayer(player);
        Service.gI().chatJustForMe(player, player.pet, "Cell Baby đã xuất hiện!");
    }

    // public void changePicPet(Player player) {
    //     byte limitPower = player.pet.nPoint.limitPower;
    //     if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
    //         player.pet.unFusion();
    //     }
    //     ChangeMapService.gI().exitMap(player.pet);
    //     player.pet.dispose();
    //     player.pet = null;
    //     createPicPet(player, limitPower);
    // }

    // public void changePicPet(Player player, int gender) {
    //     byte limitPower = player.pet.nPoint.limitPower;
    //     if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
    //         player.pet.unFusion();
    //     }
    //     ChangeMapService.gI().exitMap(player.pet);
    //     player.pet.dispose();
    //     player.pet = null;
    //     createPicPet(player, gender, limitPower);
    // }
  // public void changePicPet(Player player, int gender) {
    //     byte limitPower = player.pet.nPoint.limitPower;
    //     if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
    //         player.pet.unFusion();
    //     }
    //     ChangeMapService.gI().exitMap(player.pet);
    //     player.pet.dispose();
    //     player.pet = null;
    //     createPicPet(player, gender, limitPower);
    // }
    public void changeBlackPet(Player player) {
        byte limitPower = player.pet.nPoint.limitPower;
        if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
            player.pet.unFusion();
        }
        ChangeMapService.gI().exitMap(player.pet);
        player.pet.dispose();
        player.pet = null;
        createBlackPet(player, limitPower);
    }

    public void changeBlackPet(Player player, int gender) {
        byte limitPower = player.pet.nPoint.limitPower;
        if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
            player.pet.unFusion();
        }
        ChangeMapService.gI().exitMap(player.pet);
        player.pet.dispose();
        player.pet = null;
        createBlackPet(player, gender, limitPower);
    }

    public void changeBlackPetNow(Player player, int gender, byte... limitPower) {
        byte oldLimitPower = player.pet.nPoint.limitPower;
        if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
            player.pet.unFusion();
        }
        ChangeMapService.gI().exitMap(player.pet);
        player.pet.dispose();
        player.pet = null;
        createNewPet(player, false, false, false, true, (byte) gender);
        player.pet.nPoint.limitPower = limitPower != null && limitPower.length == 1 ? limitPower[0] : oldLimitPower;
        player.pet.nPoint.initPowerLimit();
        player.pet.joinMapMaster();
        PlayerDAO.updatePlayer(player);
        Service.gI().chatJustForMe(player, player.pet, "Fide Baby đã xuất hiện!");
    }

    public void changeNamePet(Player player, String name) {
        try {
            if (!InventoryService.gI().isExistItemBag(player, 400)) {
                Service.gI().sendThongBao(player, "Bạn cần thẻ đặt tên đệ tử, mua tại Santa");
                return;
            } else if (Util.haveSpecialCharacter(name)) {
                Service.gI().sendThongBao(player, "Tên không được chứa ký tự đặc biệt");
                return;
            } else if (name.length() > 10) {
                Service.gI().sendThongBao(player, "Tên quá dài");
                return;
            }
            ChangeMapService.gI().exitMap(player.pet);
            player.pet.name = "$" + name.toLowerCase().trim();
            InventoryService.gI().subQuantityItemsBag(player, InventoryService.gI().findItemBag(player, 400), 1);
            new Thread(() -> {
                try {
                    Thread.sleep(1000);
                    Service.gI().chatJustForMe(player, player.pet, "Cảm ơn sư phụ đã đặt cho con tên " + name);
                } catch (Exception e) {
                }
            }).start();
        } catch (Exception ex) {

        }
    }

    private int[] getDataPetNormal() {
        int[] petData = new int[5];
        petData[0] = Util.nextInt(1, 200) * 20; // hp
        petData[1] = Util.nextInt(1, 200) * 20; // mp
        petData[2] = Util.nextInt(1, 100); // dame
        petData[3] = Util.nextInt(1, 50); // def
        petData[4] = Util.nextInt(1, 10); // crit
        return petData;
    }

    private int[] getDataPetMabu() {
        int[] petData = new int[5];
        petData[0] = Util.nextInt(1, 300) * 20; // hp
        petData[1] = Util.nextInt(1, 300) * 20; // mp
        petData[2] = Util.nextInt(1, 150); // dame
        petData[3] = Util.nextInt(1, 75); // def
        petData[4] = Util.nextInt(1, 15); // crit
        return petData;
    }

    // private int[] getDataPetPic() {
    //     int[] petData = new int[5];
    //     petData[0] = Util.nextInt(40, 115) * 20; // hp
    //     petData[1] = Util.nextInt(40, 115) * 20; // mp
    //     petData[2] = Util.nextInt(70, 140); // dame
    //     petData[3] = Util.nextInt(9, 50); // def
    //     petData[4] = Util.nextInt(0, 2); // crit
    //     return petData;
    // }

    public void createPetFideNhi(Player player, boolean isChange, byte gender) {
        byte limitPower;
        if (isChange) {
            limitPower = player.pet.nPoint.limitPower;
            if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
                player.pet.unFusion();
            }
            ChangeMapService.gI().exitMap(player.pet);
            player.pet.dispose();
            player.pet = null;
        } else {
            limitPower = 1;
        }
        new Thread(() -> {
            try {
                Pet pet = new Pet(player);
                pet.name = "$Fide Nhí";
                pet.gender = gender;
                pet.id = -player.id;
                pet.nPoint.power = 1500000;
                pet.typePet = 5;
                pet.nPoint.stamina = (short) 1000;
                pet.nPoint.maxStamina = (short) 1000;
                pet.nPoint.hpg = Util.nextInt(2000, 5000);
                pet.nPoint.mpg = Util.nextInt(2000, 5000);
                pet.nPoint.hpMax = Util.nextInt(2000, 5000);
                pet.nPoint.mpMax = Util.nextInt(2000, 5000);
                pet.nPoint.dameg = Util.nextInt(200, 300);
                pet.nPoint.defg = Util.nextInt(10, 30);
                pet.nPoint.critg = 5;
                for (int i = 0; i < 7; i++) {
                    pet.inventory.itemsBody.add(ItemService.gI().createItemNull());
                }
                pet.playerSkill.skills.add(SkillUtil.createSkill(Util.nextInt(0, 2) * 2, 1));
                for (int i = 0; i < 4; i++) {
                    pet.playerSkill.skills.add(SkillUtil.createEmptySkill());
                }
                pet.nPoint.setFullHpMp();
                player.pet = pet;
                player.pet.nPoint.limitPower = limitPower;
                player.pointfusion.setHpFusion(0);
                player.pointfusion.setMpFusion(0);
                player.pointfusion.setDameFusion(0);
                Thread.sleep(1000);
                Service.gI().chatJustForMe(player, player.pet, "|0|Sư phụ ơi con đây nè!!!");
            } catch (Exception e) {
                e.printStackTrace();
            }
        }).start();
    }

    public void createPetCellNhi(Player player, boolean isChange, byte gender) {
        byte limitPower;
        if (isChange) {
            limitPower = player.pet.nPoint.limitPower;
            if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
                player.pet.unFusion();
            }
            ChangeMapService.gI().exitMap(player.pet);
            player.pet.dispose();
            player.pet = null;
        } else {
            limitPower = 1;
        }
        new Thread(() -> {
            try {
                Pet pet = new Pet(player);
                pet.name = "$Cell Nhí";
                pet.gender = gender;
                pet.id = -player.id;
                pet.nPoint.power = 1500000;
                pet.typePet = 6;
                pet.nPoint.stamina = (short) 1000;
                pet.nPoint.maxStamina = (short) 1000;
                pet.nPoint.hpg = Util.nextInt(2000, 5000);
                pet.nPoint.mpg = Util.nextInt(2000, 5000);
                pet.nPoint.hpMax = Util.nextInt(2000, 5000);
                pet.nPoint.mpMax = Util.nextInt(2000, 5000);
                pet.nPoint.dameg = Util.nextInt(200, 300);
                pet.nPoint.defg = Util.nextInt(25, 50);
                pet.nPoint.critg = 5;
                for (int i = 0; i < 7; i++) {
                    pet.inventory.itemsBody.add(ItemService.gI().createItemNull());
                }
                pet.playerSkill.skills.add(SkillUtil.createSkill(Util.nextInt(0, 2) * 2, 1));
                for (int i = 0; i < 4; i++) {
                    pet.playerSkill.skills.add(SkillUtil.createEmptySkill());
                }
                pet.nPoint.setFullHpMp();
                player.pet = pet;
                player.pet.nPoint.limitPower = limitPower;
                player.pointfusion.setHpFusion(0);
                player.pointfusion.setMpFusion(0);
                player.pointfusion.setDameFusion(0);
                Thread.sleep(1000);
                Service.gI().chatJustForMe(player, player.pet, "|0|Sư phụ ơi con đây nè!!!");
            } catch (Exception e) {
                e.printStackTrace();
            }
        }).start();
    }

    // public void createPetBuuNhi(Player player, boolean isChange, byte gender) {
    //     byte limitPower;
    //     if (isChange) {
    //         limitPower = player.pet.nPoint.limitPower;
    //         if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
    //             player.pet.unFusion();
    //         }
    //         ChangeMapService.gI().exitMap(player.pet);
    //         player.pet.dispose();
    //         player.pet = null;
    //     } else {
    //         limitPower = 1;
    //     }
    //     new Thread(() -> {
    //         try {
    //             Pet pet = new Pet(player);
    //             pet.name = "$Bưu Nhí";
    //             pet.gender = gender;
    //             pet.id = -player.id;
    //             pet.nPoint.power = 1500000;
    //             pet.typePet = 7;
    //             pet.nPoint.stamina = (short) 1000;
    //             pet.nPoint.maxStamina = (short) 1000;
    //             pet.nPoint.hpg = Util.nextInt(2000, 5000);
    //             pet.nPoint.mpg = Util.nextInt(2000, 5000);
    //             pet.nPoint.hpMax = Util.nextInt(2000, 5000);
    //             pet.nPoint.mpMax = Util.nextInt(2000, 5000);
    //             pet.nPoint.dameg = Util.nextInt(200, 300);
    //             pet.nPoint.defg = Util.nextInt(50, 100);
    //             pet.nPoint.critg = 15;
    //             for (int i = 0; i < 7; i++) {
    //                 pet.inventory.itemsBody.add(ItemService.gI().createItemNull());
    //             }
    //             pet.playerSkill.skills.add(SkillUtil.createSkill(Util.nextInt(0, 2) * 2, 1));
    //             for (int i = 0; i < 4; i++) {
    //                 pet.playerSkill.skills.add(SkillUtil.createEmptySkill());
    //             }
    //             pet.nPoint.setFullHpMp();
    //             player.pet = pet;
    //             player.pet.nPoint.limitPower = limitPower;
    //             player.pointfusion.setHpFusion(0);
    //             player.pointfusion.setMpFusion(0);
    //             player.pointfusion.setDameFusion(0);
    //             Thread.sleep(1000);
    //             Service.gI().chatJustForMe(player, player.pet, "|0|Sư phụ ơi con đây nè!!!");
    //         } catch (Exception e) {
    //             e.printStackTrace();
    //         }
    //     }).start();
    // }

    // public void createPetBerrusNhi(Player player, boolean isChange, byte gender) {
    //     byte limitPower;
    //     if (isChange) {
    //         limitPower = player.pet.nPoint.limitPower;
    //         if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
    //             player.pet.unFusion();
    //         }
    //         ChangeMapService.gI().exitMap(player.pet);
    //         player.pet.dispose();
    //         player.pet = null;
    //     } else {
    //         limitPower = 1;
    //     }
    //     new Thread(() -> {
    //         try {
    //             Pet pet = new Pet(player);
    //             pet.name = "$Berrus nhí";
    //             pet.gender = gender;
    //             pet.id = -player.id;
    //             pet.nPoint.power = 1500000;
    //             pet.typePet = 8;
    //             pet.nPoint.stamina = (short) 1000;
    //             pet.nPoint.maxStamina = (short) 1000;
    //             pet.nPoint.hpg = Util.nextInt(2000, 5000);
    //             pet.nPoint.mpg = Util.nextInt(2000, 5000);
    //             pet.nPoint.hpMax = Util.nextInt(2000, 5000);
    //             pet.nPoint.mpMax = Util.nextInt(2000, 5000);
    //             pet.nPoint.dameg = Util.nextInt(200, 300);
    //             pet.nPoint.defg = Util.nextInt(50, 100);
    //             pet.nPoint.critg = 15;
    //             for (int i = 0; i < 7; i++) {
    //                 pet.inventory.itemsBody.add(ItemService.gI().createItemNull());
    //             }
    //             pet.playerSkill.skills.add(SkillUtil.createSkill(Util.nextInt(0, 2) * 2, 1));
    //             for (int i = 0; i < 4; i++) {
    //                 pet.playerSkill.skills.add(SkillUtil.createEmptySkill());
    //             }
    //             pet.nPoint.setFullHpMp();
    //             player.pet = pet;
    //             player.pet.nPoint.limitPower = limitPower;
    //             player.pointfusion.setHpFusion(0);
    //             player.pointfusion.setMpFusion(0);
    //             player.pointfusion.setDameFusion(0);
    //             Thread.sleep(1000);
    //             Service.gI().chatJustForMe(player, player.pet, "|0|Sư phụ ơi con đây nè!!!");
    //         } catch (Exception e) {
    //             e.printStackTrace();
    //         }
    //     }).start();
    // }

    // public void createPetAdrBeach(Player player, boolean isChange, byte gender) {
    //     byte limitPower;
    //     if (isChange) {
    //         limitPower = player.pet.nPoint.limitPower;
    //         if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
    //             player.pet.unFusion();
    //         }
    //         ChangeMapService.gI().exitMap(player.pet);
    //         player.pet.dispose();
    //         player.pet = null;
    //     } else {
    //         limitPower = 1;
    //     }
    //     new Thread(() -> {
    //         try {
    //             Pet pet = new Pet(player);
    //             pet.name = "$Adr Bãi biển";
    //             pet.gender = gender;
    //             pet.id = -player.id;
    //             pet.nPoint.power = 1500000;
    //             pet.typePet = 9;
    //             pet.nPoint.stamina = (short) 1000;
    //             pet.nPoint.maxStamina = (short) 1000;
    //             pet.nPoint.hpg = Util.nextInt(2000, 5000);
    //             pet.nPoint.mpg = Util.nextInt(2000, 5000);
    //             pet.nPoint.hpMax = Util.nextInt(2000, 5000);
    //             pet.nPoint.mpMax = Util.nextInt(2000, 5000);
    //             pet.nPoint.dameg = Util.nextInt(200, 300);
    //             pet.nPoint.defg = Util.nextInt(50, 100);
    //             pet.nPoint.critg = 15;
    //             for (int i = 0; i < 7; i++) {
    //                 pet.inventory.itemsBody.add(ItemService.gI().createItemNull());
    //             }
    //             pet.playerSkill.skills.add(SkillUtil.createSkill(Util.nextInt(0, 2) * 2, 1));
    //             for (int i = 0; i < 4; i++) {
    //                 pet.playerSkill.skills.add(SkillUtil.createEmptySkill());
    //             }
    //             pet.nPoint.setFullHpMp();
    //             player.pet = pet;
    //             player.pet.nPoint.limitPower = limitPower;
    //             player.pointfusion.setHpFusion(0);
    //             player.pointfusion.setMpFusion(0);
    //             player.pointfusion.setDameFusion(0);
    //             Thread.sleep(1000);
    //             Service.gI().chatJustForMe(player, player.pet, "|0|Sư phụ ơi con đây nè!!!");
    //         } catch (Exception e) {
    //             e.printStackTrace();
    //         }
    //     }).start();
    // }

    // public void createPetMabuGay(Player player, boolean isChange, byte gender) {
    //     byte limitPower;
    //     if (isChange) {
    //         limitPower = player.pet.nPoint.limitPower;
    //         if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
    //             player.pet.unFusion();
    //         }
    //         ChangeMapService.gI().exitMap(player.pet);
    //         player.pet.dispose();
    //         player.pet = null;
    //     } else {
    //         limitPower = 1;
    //     }
    //     new Thread(() -> {
    //         try {
    //             Pet pet = new Pet(player);
    //             pet.name = "$Mabu gầy";
    //             pet.gender = gender;
    //             pet.id = -player.id;
    //             pet.nPoint.power = 1500000;
    //             pet.typePet = 10;
    //             pet.nPoint.stamina = (short) 1000;
    //             pet.nPoint.maxStamina = (short) 1000;
    //             pet.nPoint.hpg = Util.nextInt(2000, 5000);
    //             pet.nPoint.mpg = Util.nextInt(2000, 5000);
    //             pet.nPoint.hpMax = Util.nextInt(2000, 5000);
    //             pet.nPoint.mpMax = Util.nextInt(2000, 5000);
    //             pet.nPoint.dameg = Util.nextInt(200, 300);
    //             pet.nPoint.defg = Util.nextInt(50, 100);
    //             pet.nPoint.critg = 15;
    //             for (int i = 0; i < 7; i++) {
    //                 pet.inventory.itemsBody.add(ItemService.gI().createItemNull());
    //             }
    //             pet.playerSkill.skills.add(SkillUtil.createSkill(Util.nextInt(0, 2) * 2, 1));
    //             for (int i = 0; i < 4; i++) {
    //                 pet.playerSkill.skills.add(SkillUtil.createEmptySkill());
    //             }
    //             pet.nPoint.setFullHpMp();
    //             player.pet = pet;
    //             player.pet.nPoint.limitPower = limitPower;
    //             player.pointfusion.setHpFusion(0);
    //             player.pointfusion.setMpFusion(0);
    //             player.pointfusion.setDameFusion(0);
    //             Thread.sleep(1000);
    //             Service.gI().chatJustForMe(player, player.pet, "|0|Sư phụ ơi con đây nè!!!");
    //         } catch (Exception e) {
    //             e.printStackTrace();
    //         }
    //     }).start();
    // }

    private void createNewPet(Player player, boolean isMabu, boolean isBeerus, boolean isPic, boolean isBlack,
            byte... gender) {
        int[] data = isMabu ? getDataPetMabu() : getDataPetNormal();
        Pet pet = new Pet(player);
        pet.name = "$" + (isMabu ? "Mabư" : isBeerus ? "Cell Baby" : isBlack ? "Fide Baby" : "Đệ tử");
        //isBeerus ? "Beerus" : isPic ? "Pic" : 
        // isBlack ? "Black" : 
        pet.gender = (gender != null && gender.length != 0) ? gender[0] : (byte) Util.nextInt(0, 2);
        pet.id = player.isPl() ? -player.id : -Math.abs(player.id) - 100000;
   if (isBlack) {
    pet.nPoint.power = 1500000;
} else if (isMabu) {
    pet.nPoint.power = 1500000;
} else if (isBeerus) {
    pet.nPoint.power = 1500000;
} else if (isPic) {
    pet.nPoint.power = 500000; // ví dụ
} else {
    pet.nPoint.power = 2000;
}
        pet.typePet = (byte) (isMabu ? 1 : isBeerus ? 2 : isPic ? 3 : isBlack ? 4 : 0);
        if (isBeerus) {
            pet.head = 1479;
        } else if (isBlack) {
            pet.head = 1476;
        }
        pet.nPoint.stamina = 1000;
        pet.nPoint.maxStamina = 1000;
        pet.nPoint.hpg = data[0];
        pet.nPoint.mpg = data[1];
        pet.nPoint.hpMax = data[0];
        pet.nPoint.mpMax = data[1];
        pet.nPoint.dameg = data[2];
        pet.nPoint.defg = data[3];
        pet.nPoint.critg = data[4];
        for (int i = 0; i < 8; i++) {
            pet.inventory.itemsBody.add(ItemService.gI().createItemNull());
        }
            Skill skill = null;
int tiLeKame = 50;      // 80%
int tiLeMasenko = 30;   // 10%
int tiLeAntomic = 20;   // 10%

int rd = Util.nextInt(1, 100); // random từ 1-100
if (rd <= tiLeKame) {
    skill = SkillUtil.createSkill(Skill.DEMON, 1);
} else if (rd <= tiLeKame + tiLeMasenko) {
    skill = SkillUtil.createSkill(Skill.DEMON, 1);
} else if (rd <= tiLeKame + tiLeMasenko + tiLeAntomic) {
    skill = SkillUtil.createSkill(Skill.GALICK, 1);
}

// Set cooldown nếu muốn
if (skill != null) skill.coolDown = 500;

// Thêm skill này vào pet, các slot còn lại là rỗng
pet.playerSkill.skills.add(skill);             // skill chính
for (int i = 0; i < 3; i++) {
    pet.playerSkill.skills.add(SkillUtil.createEmptySkill());
}
        pet.nPoint.setFullHpMp();
        player.pet = pet;
    }

    public void createNormalPetSuperGender(Player player, int gender, byte type) {
        new Thread(() -> {
            try {
                createNewPetSuperGender(player, (byte) gender, type);
                Thread.sleep(1000);
                Service.gI().chatJustForMe(player, player.pet, "Ta là đệ tử mới của ngươi!");
            } catch (Exception e) {
                e.printStackTrace();
            }
        }).start();
    }

    public void changeNormalPetSuperGenderNow(Player player, int gender, byte type) {
        deletePet(player);
        createNewPetSuperGender(player, (byte) gender, type);
        player.pet.nPoint.initPowerLimit();
        player.pet.joinMapMaster();
        PlayerDAO.updatePlayer(player);
        Service.gI().chatJustForMe(player, player.pet, "Ta là đệ tử mới của ngươi!");
    }

    public void createNormalPetSuper(Player player, int gender, byte type) {
        new Thread(() -> {
            try {
                createNewPetSuper(player, (byte) gender, type);
                Thread.sleep(1000);
                Service.gI().chatJustForMe(player, player.pet, "Ta là đệ tử mới của ngươi!");
            } catch (Exception e) {
                e.printStackTrace();
            }
        }).start();
    }

    private void createNewPetSuper(Player player, byte gender, byte type) {
        int[] data = getDataPetNormal();
        Pet pet = new Pet(player);
        pet.name = "$" + "Đệ tử";
        pet.gender = (byte) Util.nextInt(0, 2);
        pet.id = -player.id;
        pet.nPoint.power = 2000;
        pet.typePet = type;
        pet.nPoint.stamina = 1000;
        pet.nPoint.maxStamina = 1000;
        pet.nPoint.hpg = data[0];
        pet.nPoint.mpg = data[1];
        pet.nPoint.hpMax = data[0];
        pet.nPoint.mpMax = data[1];
        pet.nPoint.dameg = data[2];
        pet.nPoint.defg = data[3];
        pet.nPoint.critg = data[4];
        for (int i = 0; i < 8; i++) {
            pet.inventory.itemsBody.add(ItemService.gI().createItemNull());
        }
       Skill skill = null;
int tiLeKame = 50;      // 80%
int tiLeMasenko = 30;   // 10%
int tiLeAntomic = 20;   // 10%

int rd = Util.nextInt(1, 100); // random từ 1-100
if (rd <= tiLeKame) {
    skill = SkillUtil.createSkill(Skill.DEMON, 1);
} else if (rd <= tiLeKame + tiLeMasenko) {
    skill = SkillUtil.createSkill(Skill.DEMON, 1);
} else if (rd <= tiLeKame + tiLeMasenko + tiLeAntomic) {
    skill = SkillUtil.createSkill(Skill.GALICK, 1);
}

// Set cooldown nếu muốn
if (skill != null) skill.coolDown = 500;

// Thêm skill này vào pet, các slot còn lại là rỗng
pet.playerSkill.skills.add(skill);             // skill chính
for (int i = 0; i < 3; i++) {
    pet.playerSkill.skills.add(SkillUtil.createEmptySkill());
}
        pet.nPoint.setFullHpMp();
        player.pet = pet;
        player.pointfusion.setHpFusion(0);
        player.pointfusion.setMpFusion(0);
        player.pointfusion.setDameFusion(0);
    }

    private void createNewPetSuperGender(Player player, byte gender, byte type) {
        int[] data = getDataPetNormal();
        Pet pet = new Pet(player);
        pet.name = "$" + (type == 1 ? "[Broly]Mabư"
                : type == 2 ? "[Broly]Beerus" : type == 3 ? "[Broly]Pic" : type == 4 ? "[Broly]Black" : "Đệ tử");
        pet.gender = gender;
        pet.id = -player.id;
        pet.nPoint.power = 2000;
        pet.typePet = type;
        pet.nPoint.stamina = 1000;
        pet.nPoint.maxStamina = 1000;
        pet.nPoint.hpg = data[0];
        pet.nPoint.mpg = data[1];
        pet.nPoint.hpMax = data[0];
        pet.nPoint.mpMax = data[1];

        pet.nPoint.dameg = data[2];
        pet.nPoint.defg = data[3];
        pet.nPoint.critg = data[4];
        for (int i = 0; i < 8; i++) {
            pet.inventory.itemsBody.add(ItemService.gI().createItemNull());
        }
              Skill skill = null;
int tiLeKame = 50;      // 80%
int tiLeMasenko = 30;   // 10%
int tiLeAntomic = 20;   // 10%

int rd = Util.nextInt(1, 100); // random từ 1-100
if (rd <= tiLeKame) {
    skill = SkillUtil.createSkill(Skill.DEMON, 1);
} else if (rd <= tiLeKame + tiLeMasenko) {
    skill = SkillUtil.createSkill(Skill.DEMON, 1);
} else if (rd <= tiLeKame + tiLeMasenko + tiLeAntomic) {
    skill = SkillUtil.createSkill(Skill.GALICK, 1);
}

// Set cooldown nếu muốn
if (skill != null) skill.coolDown = 500;

// Thêm skill này vào pet, các slot còn lại là rỗng
pet.playerSkill.skills.add(skill);             // skill chính
for (int i = 0; i < 3; i++) {
    pet.playerSkill.skills.add(SkillUtil.createEmptySkill());

        }
        pet.nPoint.setFullHpMp();
        player.pet = pet;
        player.pointfusion.setHpFusion(0);
        player.pointfusion.setMpFusion(0);
        player.pointfusion.setDameFusion(0);
    }

    public static void Pet2(Player pl, int h, int b, int l) {
        if (pl.newPet != null) {
            pl.newPet.dispose();
        }
        pl.newPet = new NewPet(pl, (short) h, (short) b, (short) l);
        pl.newPet.name = "$";
        pl.newPet.gender = pl.gender;
        pl.newPet.nPoint.tiemNang = 1;
        pl.newPet.nPoint.power = 1;
        pl.newPet.nPoint.limitPower = 1;
        pl.newPet.nPoint.hpg = 500000;
        pl.newPet.nPoint.mpg = 500000;
        pl.newPet.nPoint.hp = 500000;
        pl.newPet.nPoint.mp = 500000;
        pl.newPet.nPoint.dameg = 1;
        pl.newPet.nPoint.defg = 1;
        pl.newPet.nPoint.critg = 1;
        pl.newPet.nPoint.stamina = 1;
        pl.newPet.nPoint.setBasePoint();
        pl.newPet.nPoint.setFullHpMp();
    }

    public void deletePet(Player player) {
        Pet pet = player.pet;
        if (pet != null) {
            if (player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
                pet.unFusion();
            }
            ChangeMapService.gI().exitMap(pet);
            pet.dispose();
            player.pet = null;
        }
    }
}
