package server;

import GiaCatLuongManager.GiaCatLuongManager;
import consts.cn;
import data.ItemData;
import item.Item;
import item.Item.ItemOption;
import javax.swing.*;
import java.awt.*;
import java.awt.event.ActionEvent;
import java.awt.event.MouseAdapter;
import java.awt.event.MouseEvent;
import java.awt.event.WindowAdapter;
import java.awt.event.WindowEvent;
import java.io.BufferedReader;
import java.io.BufferedWriter;
import java.io.FileWriter;
import java.io.IOException;
import java.lang.management.ManagementFactory;
import com.sun.management.OperatingSystemMXBean;
import java.time.LocalTime;
import java.util.List;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.prefs.Preferences;
import jdbc.daos.EventDAO;
import jdbc.daos.NDVSqlFetcher;
import jdbc.daos.PlayerDAO;
import models.Consign.ConsignShopManager;
import network.SessionManager;
import services.ClanService;
import utils.Logger;

import player.Player;
import services.ItemService;

import services.Service;

public class ServerManagerUI extends JFrame {

    private JLabel threadCountLabel;
    private JLabel cpuLabel;
    private JLabel ramLabel;
    private JLabel messageLabel;
    private JLabel countdownLabel;
    private JLabel info;

    public ServerManagerUI() {
        Preferences.userNodeForPackage(ServerManagerUI.class);
        setTitle("Ngoc Rong Kame - May Chu " + cn.SV + "");
        ImageIcon icon = new ImageIcon(getClass().getResource("icon.png"));
        setIconImage(icon.getImage());
        setSize(400, 230); // Tăng kích thước cửa sổ để chứa thêm các nút
        setResizable(false); // Khóa chỉnh kích thước
        setDefaultCloseOperation(JFrame.DO_NOTHING_ON_CLOSE);
        setLocationRelativeTo(null);
        addWindowListener(new WindowAdapter() {
            @Override
            public void windowClosing(WindowEvent e) {
                confirmExit();
            }
        });

        initUI();
        setVisible(true);

        ServerManager.gI().run();
        GiaCatLuongManager.getInstance().startAutoSave();
    }

    private void initUI() {
        JPanel panel = new JPanel();
        panel.setBackground(new Color(240, 240, 240));
        panel.setLayout(new BoxLayout(panel, BoxLayout.Y_AXIS));
        getContentPane().add(panel);

        JButton maintenanceButton = createStyledButton("Bảo trì", new Color(220, 53, 69), new Color(255, 165, 0));
        maintenanceButton.setAlignmentX(Component.CENTER_ALIGNMENT);
        maintenanceButton.addActionListener(e -> showMaintenanceDialog());
        panel.add(Box.createVerticalStrut(10));
        panel.add(maintenanceButton);
        panel.add(Box.createVerticalStrut(12));

        messageLabel = createCenteredLabel("");
        countdownLabel = createCenteredLabel("");
        info = createCenteredLabel("");
        threadCountLabel = createCenteredLabel("Số thread: 0");
        cpuLabel = createCenteredLabel("CPU game: 0.0% / VPS: 0.0%");
        ramLabel = createCenteredLabel("RAM game: 0.0% / VPS: 0.0%");

        panel.add(messageLabel);
        panel.add(countdownLabel);
        panel.add(info);
        panel.add(threadCountLabel);
        panel.add(Box.createVerticalStrut(6));
        panel.add(cpuLabel);
        panel.add(Box.createVerticalStrut(6));
        panel.add(ramLabel);

        startScheduledTasks();
    }

    private JLabel createCenteredLabel(String text) {
        JLabel label = new JLabel(text);
        label.setAlignmentX(Component.CENTER_ALIGNMENT);
        return label;
    }

    private JButton createStyledButton(String text, Color normalColor, Color hoverColor) {
        JButton button = new JButton(text) {
            @Override
            protected void paintComponent(Graphics g) {
                // Vẽ nút bo tròn
                Graphics2D g2 = (Graphics2D) g.create();
                g2.setRenderingHint(RenderingHints.KEY_ANTIALIASING, RenderingHints.VALUE_ANTIALIAS_ON);
                g2.setColor(getModel().isRollover() ? hoverColor : normalColor); // Đổi màu khi hover
                g2.fillRoundRect(0, 0, getWidth(), getHeight(), 20, 20); // Bo tròn góc 20px
                g2.setColor(Color.WHITE);
                g2.drawString(text, getWidth() / 2 - g2.getFontMetrics().stringWidth(text) / 2, getHeight() / 2 + 5);
                g2.dispose();
            }

            @Override
            protected void paintBorder(Graphics g) {
                // Không vẽ viền mặc định
            }
        };
        button.setContentAreaFilled(false); // Không tô nền mặc định
        button.setFocusPainted(false); // Loại bỏ viền focus
        button.setFont(new Font("Arial", Font.BOLD, 12)); // Font chữ đậm
        button.setForeground(Color.WHITE); // Màu chữ trắng
        button.setPreferredSize(new Dimension(120, 40));
        button.setMaximumSize(new Dimension(120, 40)); // Kích thước nút

        // Thêm hiệu ứng hover
        button.addMouseListener(new MouseAdapter() {
            @Override
            public void mouseEntered(MouseEvent e) {
                button.setBackground(hoverColor); // Đổi màu khi rê chuột vào
                button.repaint();
            }

            @Override
            public void mouseExited(MouseEvent e) {
                button.setBackground(normalColor); // Trở lại màu ban đầu khi chuột rời khỏi
                button.repaint();
            }
        });

        return button;
    }

    private void startScheduledTasks() {
        ScheduledExecutorService statsExecutor = Executors.newSingleThreadScheduledExecutor();
        statsExecutor.scheduleAtFixedRate(() -> {
            int threadCount = Thread.activeCount();
            int onlineCount = Client.gI().getPlayers().size();
            int sessionCount = SessionManager.gI().getSessions().size();
            OperatingSystemMXBean os = (OperatingSystemMXBean) ManagementFactory.getOperatingSystemMXBean();
            double processCpu = Math.max(0D, os.getProcessCpuLoad()) * 100D;
            double systemCpu = Math.max(0D, os.getSystemCpuLoad()) * 100D;
            long gameRam = Runtime.getRuntime().totalMemory() - Runtime.getRuntime().freeMemory();
            long totalRam = os.getTotalPhysicalMemorySize();
            long usedVpsRam = totalRam - os.getFreePhysicalMemorySize();
            double gameRamPercent = totalRam > 0 ? gameRam * 100D / totalRam : 0D;
            double vpsRamPercent = totalRam > 0 ? usedVpsRam * 100D / totalRam : 0D;

            SwingUtilities.invokeLater(() -> {
                threadCountLabel.setText("Số thread: " + threadCount + "  Online: " + onlineCount + "  Session: " + sessionCount);
                cpuLabel.setText(String.format("CPU game: %.1f%% / VPS: %.1f%%", processCpu, systemCpu));
                ramLabel.setText(String.format("RAM game: %.1f%% / VPS: %.1f%%", gameRamPercent, vpsRamPercent));
            });
        }, 1, 1, TimeUnit.SECONDS);
    }

    private String formatMemory(long bytes) {
        double gigabytes = bytes / (1024D * 1024D * 1024D);
        if (gigabytes >= 1D) {
            return String.format("%.2f GB", gigabytes);
        }
        return String.format("%.0f MB", bytes / (1024D * 1024D));
    }

    private void saveData(ActionEvent e) {
        Logger.success("Đang tiến hành lưu data\n");
        network.server.GiaCatLuongServer.gI().stopConnect();

        Maintenance.isRunning = false;
        try {
            Logger.error("Đang tiến hành lưu data bang hội\n");
            ClanService.gI().close();
            Thread.sleep(1000);
            Logger.success("Lưu dữ liệu bang hội thành công\n");
        } catch (Exception ex) {
            Logger.error("Lỗi lưu dữ liệu bang hội\n");
        }
        try {
            Logger.error("Đang tiến hành lưu data ký gửi\n");
            ConsignShopManager.gI().save();
            Thread.sleep(1000);
            Logger.success("Lưu dữ liệu ký gửi thành công\n");
        } catch (Exception ex) {
            Logger.error("Lỗi lưu dữ liệu ký gửi\n");
        }

        try {
            Logger.error("Đang tiến hành đẩy người chơi\n");
            Client.gI().close();
            EventDAO.save();
            Thread.sleep(1000);
            Logger.success("Lưu dữ liệu người dùng thành công\n");
        } catch (Exception ex) {

            Logger.error("Lỗi lưu dữ liệu người dùng\n");
        }
        System.exit(0);
    }

    private void showMaintenanceDialog() {
        int dialogButton = JOptionPane.YES_NO_OPTION;
        int dialogResult = JOptionPane.showConfirmDialog(this, "Bắt đầu bảo trì?", "Bảo trì", dialogButton);
        if (dialogResult == 0) {
            Logger.error("Server tiến hành bảo trì");
            Maintenance.gI().start(15);
        }
    }

    private void scheduleMaintenance(JComboBox<Integer> hoursComboBox, JComboBox<Integer> minutesComboBox,
            JComboBox<Integer> secondsComboBox) {
        int hours = hoursComboBox.getItemAt(hoursComboBox.getSelectedIndex());
        int minutes = minutesComboBox.getItemAt(minutesComboBox.getSelectedIndex());
        int seconds = secondsComboBox.getItemAt(secondsComboBox.getSelectedIndex());
        if (hours == -1 || minutes == -1 || seconds == -1) {
            JOptionPane.showMessageDialog(this, "Chạy server không cần hẹn bảo trì?");
            return;
        }

        try (BufferedWriter writer = new BufferedWriter(new FileWriter("maintenanceConfig.txt"))) {
            writer.write(hours + "\n");
            writer.write(minutes + "\n");
            writer.write(seconds + "\n");
            writer.flush();
        } catch (IOException e) {
            e.printStackTrace();
        }

        AtomicBoolean timeReached = new AtomicBoolean(false);
        info.setText("Đã cài đặt quá trình bảo trì tự động vào lúc " + hours + ":" + minutes + ":" + seconds);
        new Thread(() -> {
            while (!timeReached.get()) {
                try {
                    LocalTime currentTime = LocalTime.now();
                    if (hours == currentTime.getHour() && minutes == currentTime.getMinute()
                            && seconds == currentTime.getSecond()) {
                        performMaintenance();
                        timeReached.set(true);
                    }
                    Thread.sleep(1000);
                } catch (Exception e) {

                }
            }
        }).start();
    }

    private void performMaintenance() {
        Maintenance.gI().start(15);
    }

    private void confirmExit() {
        int dialogButton = JOptionPane.YES_NO_OPTION;
        int dialogResult = JOptionPane.showConfirmDialog(this, "Bạn có chắc muốn thoát chương trình?", "Thoát",
                dialogButton);
        if (dialogResult == 0) {
            System.exit(0);
        }
    }

    // Thêm phương thức Buff VND
    public void createFormBuffVND() {
        JPanel panel = new JPanel(new GridLayout(3, 2));
        JTextField nameField = new JTextField();
        JTextField vndField = new JTextField();

        panel.add(new JLabel("Tên người chơi:"));
        panel.add(nameField);
        panel.add(new JLabel("Cá cần buff:"));
        panel.add(vndField);

        int result = JOptionPane.showConfirmDialog(this, panel, "Buff Cá", JOptionPane.OK_CANCEL_OPTION,
                JOptionPane.PLAIN_MESSAGE);
        if (result == JOptionPane.OK_OPTION) {
            try {
                String name = nameField.getText().trim();
                int addcash = Integer.parseInt(vndField.getText().trim());

                // Kiểm tra người chơi có tồn tại không
                Player player = NDVSqlFetcher.loadPlayerByName(name);
                if (player != null) {
                    // Nếu người chơi tồn tại, thực hiện buff VND
                    buffVND(name, addcash);
                } else {
                    JOptionPane.showMessageDialog(this, "Người chơi không tồn tại", "Lỗi", JOptionPane.ERROR_MESSAGE);
                }
            } catch (NumberFormatException e) {
                JOptionPane.showMessageDialog(this, "Vui lòng nhập số hợp lệ", "Lỗi", JOptionPane.ERROR_MESSAGE);
            }
        }
    }

    private void buffVND(String name, int addcash) {
        try {
            if (PlayerDAO.addcash(name, addcash)) {
                JOptionPane.showMessageDialog(this, "Bạn đã buff cho " + name + " " + addcash + " Cá");
                Player player = Client.gI().getPlayerByName(name);
                if (player != null) {
                    player.getSession().cash += addcash;
                    Service.gI().sendThongBao(player, "Bạn vừa được cộng " + addcash + " Cá bởi Admin");
                }
            }
        } catch (Exception e) {
            e.printStackTrace();
            JOptionPane.showMessageDialog(this, "Đã có lỗi xảy ra", "Lỗi", JOptionPane.ERROR_MESSAGE);
        }
    }

    // Thêm phương thức Buff Hộp Thư
    public void createFromMailBox() {
        JPanel panel = new JPanel(new GridLayout(4, 2));
        JTextField playerNameField = new JTextField();
        JTextField itemIdField = new JTextField();
        JTextField optionField = new JTextField();
        JTextField quantityField = new JTextField();

        panel.add(new JLabel("Tên người chơi:"));
        panel.add(playerNameField);
        panel.add(new JLabel("ID Trang Bị:"));
        panel.add(itemIdField);
        panel.add(new JLabel("Chuỗi option:"));
        panel.add(optionField);
        panel.add(new JLabel("Số lượng:"));
        panel.add(quantityField);

        int result = JOptionPane.showConfirmDialog(this, panel, "Buff Hộp Thư", JOptionPane.OK_CANCEL_OPTION,
                JOptionPane.PLAIN_MESSAGE);
        if (result == JOptionPane.OK_OPTION) {
            try {
                String playerName = playerNameField.getText().trim();
                String itemIds = itemIdField.getText().trim();
                String option = optionField.getText().trim();
                int quantity = Integer.parseInt(quantityField.getText().trim());

                if (quantity > 999) {
                    JOptionPane.showMessageDialog(this, "Số lượng tối đa là 999", "Lỗi", JOptionPane.ERROR_MESSAGE);
                    return;
                }

                // Gọi phương thức xử lý gửi vật phẩm
                sendItemToMailBox(playerName, itemIds, option, quantity);
            } catch (NumberFormatException e) {
                JOptionPane.showMessageDialog(this, "Vui lòng nhập số hợp lệ", "Lỗi", JOptionPane.ERROR_MESSAGE);
            }
        }
    }

    private void sendItemToMailBox(String playerName, String itemIds, String option, int quantity) {
        new Thread(() -> {
            if (playerName.equalsIgnoreCase("all")) {
                List<Player> allPlayers = NDVSqlFetcher.getAllPlayer();
                for (Player pBuffItem : allPlayers) {
                    if (pBuffItem != null) {
                        processItemBuff(pBuffItem, itemIds, option, quantity);
                    }
                }
            } else {
                Player pBuffItem = NDVSqlFetcher.loadPlayerByName(playerName);
                if (pBuffItem != null) {
                    processItemBuff(pBuffItem, itemIds, option, quantity);
                } else {
                    JOptionPane.showMessageDialog(this, "Người chơi không tồn tại", "Lỗi", JOptionPane.ERROR_MESSAGE);
                }
            }
        }).start();
    }

    private void processItemBuff(Player pBuffItem, String itemIds, String option, int quantity) {
        String[] itemIdsArray = itemIds.split(",");
        for (String itemId : itemIdsArray) {
            int idItemBuff = Integer.parseInt(itemId);
            Item itembuff = ItemService.gI().createNewItem((short) idItemBuff, quantity);

            if (option != null && !option.isEmpty()) {
                String[] options = option.split(",");
                for (String opt : options) {
                    String[] optItem = opt.split("-");
                    int optID = Integer.parseInt(optItem[0]);
                    int param = Integer.parseInt(optItem[1]);
                    itembuff.itemOptions.add(new ItemOption(optID, param));
                }
            }

            pBuffItem.inventory.itemsMailBox.add(itembuff);
            if (NDVSqlFetcher.updateMailBox(pBuffItem)) {
                JOptionPane.showMessageDialog(this, "Đã gửi " + itembuff.template.name + " đến " + pBuffItem.name,
                        "Thành công", JOptionPane.INFORMATION_MESSAGE);
            }
        }
    }

    public static void main(String[] args) {
        SwingUtilities.invokeLater(() -> new ServerManagerUI());
    }
}
