package main

import (
	"ScanElections/database"
	"ScanElections/err_msg"
	"ScanElections/utils"
	"fmt"
	"github.com/rs/xid"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"testing"
)

var displayMenuGlobalPos *core.QPoint
var displayLargeScreenType = 2
var executeMeetingId = "cgha67brm5p0qf7vf5sg"
var rootWidth, rootHeight = utils.GetRootScreenSize()

func TestDisplay(t *testing.T) {
	// 创建Qt应用程序
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// 创建主窗口
	window := widgets.NewQMainWindow(nil, 0)

	window.SetFixedSize2(rootWidth, rootHeight)
	windowCentralWidget := widgets.NewQWidget(nil, 0)
	window.SetWindowTitle("大屏展示配置")

	windowLayout := widgets.NewQHBoxLayout()

	ticketImageList := widgets.NewQListWidget(nil)
	ticketImageList.SetViewMode(widgets.QListView__IconMode)
	ticketImageList.SetIconSize(core.NewQSize2(200, 200))
	ticketImageList.SetFixedWidth(210)
	ticketImageList.SetResizeMode(widgets.QListView__Adjust)
	ticketPics, ticketPicsErr := database.FindTicketsByMeetingId(executeMeetingId)
	if ticketPicsErr != nil {
		showErrorMessageBox(nil, err_msg.FindTicketByMeetingErr)
	}
	for _, thisT := range ticketPics {
		os.WriteFile(thisT.PictureDir, thisT.PictureData, 0644)
		pixmap := gui.NewQPixmap3(thisT.PictureDir, "", 0)
		item := widgets.NewQListWidgetItem2("", ticketImageList, 0)
		item.SetIcon(gui.NewQIcon2(pixmap))
		font := item.Font()
		font.SetPointSize(9)
		item.SetFont(font)
		fm := gui.NewQFontMetrics(font)
		elidedText := fm.ElidedText(thisT.ConferenceTitleName, core.Qt__ElideRight, 180, 0)
		item.SetText(elidedText)
		item.SetTextAlignment(int(core.Qt__AlignCenter))
		item.SetData(int(core.Qt__UserRole), core.NewQVariant1(thisT.TicketId))
	}
	windowLayout.AddWidget(ticketImageList, 0, 0)

	mainWindow := widgets.NewQWidget(nil, 0)

	ticketImageList.ConnectItemClicked(func(item *widgets.QListWidgetItem) {

		ticketId := item.Data(int(core.Qt__UserRole)).ToString()
		layout := mainWindow.Layout()
		if layout != nil {
			for i := layout.Count() - 1; i >= 0; i-- {
				layout.RemoveItem(layout.ItemAt(i))
			}
		}

		//ticketId := "cgn45s3rm5p0j51tt0eg"
		// 创建QLabel并设置背景图片
		backgroundLabel := widgets.NewQLabel(mainWindow, 0)
		backGroundImage := database.FindAllAssembliesByTicketIdAndTypeNum(ticketId, 0, displayLargeScreenType)
		//找不到就默认
		if len(backGroundImage) == 0 {
			backGroundImage = append(backGroundImage, database.FindDefaultBackGround(displayLargeScreenType))
			os.WriteFile(backGroundImage[0].BackPictureDir, backGroundImage[0].BackPictureData, 0644)
			backgroundPixmap := gui.NewQPixmap3(backGroundImage[0].BackPictureDir, "", core.Qt__AutoColor)
			// 缩放图片
			scaledPixmap := backgroundPixmap.Scaled2(
				rootWidth-210,
				rootHeight,
				core.Qt__IgnoreAspectRatio,
				core.Qt__SmoothTransformation,
			)
			backgroundLabel.SetScaledContents(true)
			backgroundLabel.SetPixmap(scaledPixmap)
			// 设置QPalette对象
			palette := mainWindow.Palette()
			brush := gui.NewQBrush3(gui.NewQColor3(0, 0, 0, 0), core.Qt__TexturePattern)
			brush.SetTexture(scaledPixmap)
			palette.SetBrush(gui.QPalette__Window, brush)
			mainWindow.SetPalette(palette)
		} else {
			//找得到就判断图or纯色
			backgroundLabel.SetObjectName(backGroundImage[0].AssemblyId)
			if backGroundImage[0].BackPictureDir != "" {
				os.WriteFile(backGroundImage[0].BackPictureDir, backGroundImage[0].BackPictureData, 0644)
				backgroundPixmap := gui.NewQPixmap3(backGroundImage[0].BackPictureDir, "", core.Qt__AutoColor)
				scaledPixmap := backgroundPixmap.Scaled2(
					rootWidth-210,
					rootHeight,
					core.Qt__IgnoreAspectRatio,
					core.Qt__SmoothTransformation,
				)
				backgroundLabel.SetScaledContents(true)
				// 让图片铺满整个 QLabel
				backgroundLabel.SetPixmap(scaledPixmap)
				// 设置QPalette对象
				palette := mainWindow.Palette()
				brush := gui.NewQBrush3(gui.NewQColor3(0, 0, 0, 0), core.Qt__TexturePattern)
				brush.SetTexture(scaledPixmap)
				palette.SetBrush(gui.QPalette__Window, brush)
				mainWindow.SetPalette(palette)
			}
			if backGroundImage[0].BackColor != "" {
				style := fmt.Sprintf("background-color:%s;", backGroundImage[0].BackColor)
				backgroundLabel.SetStyleSheet(style)
			}
		}

		//遍历文本
		texts := database.FindAllAssembliesByTicketIdAndTypeNum(ticketId, 2, displayLargeScreenType)
		for _, text := range texts {
			textLabel := widgets.NewQLabel(mainWindow, 0)
			textLabel.SetText(text.TextContent)
			textLabel.SetFixedSize2(text.TextWidth, text.TextHeight)
			textLabel.SetObjectName(text.AssemblyId)
			layout := mainWindow.Layout()
			layout.AddWidget(textLabel)
			x, _ := strconv.Atoi(fmt.Sprintf("%1.0f", text.ScreenRowX*float64(rootWidth)))
			y, _ := strconv.Atoi(fmt.Sprintf("%1.0f", text.ScreenColY*float64(rootHeight)))
			textLabel.Move2(x, y)
			textColor := gui.NewQColor6(text.TextColor)
			textColor.SetAlphaF(text.Opacity)
			style := fmt.Sprintf("color: rgba(%d, %d, %d, %f);", textColor.Red(), textColor.Green(), textColor.Blue(), textColor.AlphaF())
			textLabel.SetStyleSheet(style)
			// 创建一个 QFont 对象
			font := textLabel.Font()
			font.SetPointSize(text.TextSizeInt)
			font.SetFamily(text.TextStyle)
			if text.IsTextWeight {
				font.SetWeight(75)
			} else {
				font.SetWeight(50)
			}
			font.SetItalic(text.IsTextItalic)
			font.SetUnderline(text.IsTextUnder)
			font.SetStrikeOut(text.IsTextStrikeOut)
			// 将新的字体应用到标签
			textLabel.SetFont(font)
			textLabel.Show()
			TextLabelAction(textLabel, mainWindow)
		}

		//图
		pics := database.FindAllAssembliesByTicketIdAndTypeNum(ticketId, 1, displayLargeScreenType)
		for _, pic := range pics {
			imgLabel := widgets.NewQLabel(mainWindow, 0)
			os.WriteFile(pic.PictureDir, pic.PictureData, 0644)
			newScaledPixmap := gui.NewQPixmap3(pic.PictureDir, "", core.Qt__AutoColor)
			scaleFactor := float64(pic.PictureWidth) / float64(newScaledPixmap.Width())
			if float64(newScaledPixmap.Height())*scaleFactor > float64(pic.PictureHeight) {
				scaleFactor = float64(pic.PictureHeight) / float64(newScaledPixmap.Height())
			}
			// 缩放图片
			scaledPixmap := newScaledPixmap.Scaled2(
				int(float64(newScaledPixmap.Width())*scaleFactor),
				int(float64(newScaledPixmap.Height())*scaleFactor),
				core.Qt__IgnoreAspectRatio,
				core.Qt__SmoothTransformation,
			)
			imgLabel.SetPixmap(scaledPixmap)
			palette := imgLabel.Palette()
			brush := gui.NewQBrush3(gui.NewQColor3(0, 0, 0, 0), core.Qt__TexturePattern)
			brush.SetTexture(scaledPixmap)
			palette.SetBrush(gui.QPalette__Window, brush)
			imgLabel.SetPalette(palette)
			imgLabel.SetObjectName(pic.AssemblyId)
			x, _ := strconv.Atoi(fmt.Sprintf("%1.0f", pic.ScreenRowX*float64(rootWidth)))
			y, _ := strconv.Atoi(fmt.Sprintf("%1.0f", pic.ScreenColY*float64(rootHeight)))
			imgLabel.Move2(x, y)
			ImageLabelAction(imgLabel, scaledPixmap)
		}

		// 创建右键菜单
		menu := widgets.NewQMenu(nil)
		addTextAction := menu.AddAction("添加文本")
		addImageAction := menu.AddAction("添加图片")
		addBackImageAction := menu.AddAction("修改背景图片")
		addBackColorAction := menu.AddAction("选择纯色背景")
		returnBackAction := menu.AddAction("退出")

		// 将右键菜单连接到窗口部件的右键单击事件
		mainWindow.ConnectContextMenuEvent(func(event *gui.QContextMenuEvent) {
			pos := event.Pos()
			displayMenuGlobalPos = mainWindow.MapToGlobal(pos)
			menu.Exec2(event.GlobalPos(), nil)
		})

		// 添加文本
		addTextAction.ConnectTriggered(func(bool) {
			textLabel := widgets.NewQLabel(mainWindow, 0)
			textLabel.SetText("请输入文本")
			textLabel.SetMinimumSize2(200, 50)
			layout := mainWindow.Layout()
			layout.AddWidget(textLabel)
			textLabel.Move2(displayMenuGlobalPos.X(), displayMenuGlobalPos.Y())
			textLabel.SetObjectName(xid.New().String())
			var lc database.LargeScreenConfiguration
			lc.MeetingId = executeMeetingId
			lc.TicketId = ticketId
			lc.AssemblyId = textLabel.ObjectName()
			lc.TextContent = textLabel.Text()
			lc.ScreenRowX = float64(displayMenuGlobalPos.X()-210) / float64(rootWidth)
			lc.ScreenColY = float64(displayMenuGlobalPos.Y()) / float64(rootHeight)
			lc.LargeScreenId = displayLargeScreenType
			lc.TypeNum = 2
			qFont := textLabel.Font()
			lc.TextSizeInt = qFont.PointSize()
			lc.TextStyle = qFont.Family()
			lc.TextWidth = textLabel.Width()
			lc.TextHeight = textLabel.Height()
			//var isWeight bool
			//if qFont.Weight() == 75 {
			//	isWeight = true
			//} else {
			//	isWeight = false
			//}
			//lc.IsTextWeight = isWeight
			//lc.IsTextItalic = qFont.Italic()
			//lc.IsTextUnder = qFont.Underline()
			//lc.IsTextStrikeOut = qFont.StrikeOut()
			err := database.CreateOrUpdateLargeScreenConfigurationByAssemblyId(lc)
			if err != nil {
				showErrorMessageBox(nil, err_msg.OperateLargeScreenTextErr)
			}
			textLabel.Show()
			TextLabelAction(textLabel, mainWindow)
		})

		// 添加图像
		addImageAction.ConnectTriggered(func(bool) {
			fileDialog := widgets.NewQFileDialog(nil, 0)
			fileDialog.SetAcceptMode(widgets.QFileDialog__AcceptOpen)
			fileDialog.SetFileMode(widgets.QFileDialog__ExistingFile)
			fileDialog.SetNameFilter("图片文件 (*.jpg *.jpeg *.png *.bmp *.gif)")
			// 显示文件选择对话框，并在用户选择文件后进行处理
			fileDialog.ConnectFileSelected(func(imgFiles string) {
				var lC database.LargeScreenConfiguration
				lC.MeetingId = executeMeetingId
				lC.TicketId = ticketId
				lC.TypeNum = 1
				aId := xid.New().String()
				// 读取选择的文件内容
				if imgFiles != "" {
					file, imgErr := ioutil.ReadFile(imgFiles)
					if imgErr != nil {
						showErrorMessageBox(nil, err_msg.ImageReadErr)
						return
					}
					lC.PictureDir = imgFiles
					lC.PictureData = file
				} else {
					showErrorMessageBox(nil, err_msg.OperateLargeScreenImageErr)
					return
				}
				pixmap := gui.NewQPixmap3(lC.PictureDir, "", core.Qt__AutoColor)
				// 创建 QLabel 控件QPixmap
				imgLabel := widgets.NewQLabel(mainWindow, 0)
				imgLabel.SetPixmap(pixmap)
				imgLabel.SetFixedSize2(pixmap.Width(), pixmap.Height())
				imgLabel.Move2(displayMenuGlobalPos.X()-210, displayMenuGlobalPos.Y())
				imgLabel.SetObjectName(aId)
				lC.LargeScreenId = displayLargeScreenType
				lC.AssemblyId = aId
				lC.PictureWidth = pixmap.Width()
				lC.PictureHeight = pixmap.Height()
				lC.ScreenRowX = float64(displayMenuGlobalPos.X()-210) / float64(rootWidth)
				lC.ScreenColY = float64(displayMenuGlobalPos.Y()) / float64(rootHeight)
				imgLabel.Show()
				if err := database.CreateOrUpdateLargeScreenConfigurationByAssemblyId(lC); err != nil {
					showErrorMessageBox(nil, err_msg.ConfigurationSaveErr)
					return
				}
				showSuccessMessageBox(nil, "保存成功")
				ImageLabelAction(imgLabel, pixmap)
			})
			fileDialog.Exec()
		})

		//背景图
		addBackImageAction.ConnectTriggered(func(bool) {
			fileDialog := widgets.NewQFileDialog(nil, 0)
			fileDialog.SetAcceptMode(widgets.QFileDialog__AcceptOpen)
			fileDialog.SetFileMode(widgets.QFileDialog__ExistingFile)
			fileDialog.SetNameFilter("图片文件 (*.jpg *.jpeg *.png *.bmp *.gif)")

			fileDialog.ConnectFileSelected(func(imgFiles string) {
				var lC database.LargeScreenConfiguration
				lC.MeetingId = executeMeetingId
				lC.TicketId = ticketId
				lC.TypeNum = 0
				lC.LargeScreenId = displayLargeScreenType
				lC.LargeScreenId = displayLargeScreenType
				aId := backgroundLabel.ObjectName()
				if aId == "" {
					aId = xid.New().String()
					backgroundLabel.SetObjectName(aId)
				}
				lC.AssemblyId = aId
				if imgFiles != "" {
					file, imgErr := ioutil.ReadFile(imgFiles)
					if imgErr != nil {
						showErrorMessageBox(nil, err_msg.ImageReadErr)
						return
					}
					lC.BackPictureDir = imgFiles
					lC.BackPictureData = file
				} else {
					ground := database.FindDefaultBackGround(displayLargeScreenType)
					lC.BackPictureDir = ground.BackPictureDir
					lC.BackPictureData = ground.BackPictureData
				}
				newScaledPixmap := gui.NewQPixmap3(lC.BackPictureDir, "", core.Qt__AutoColor)
				if err := database.CreateOrUpdateLargeScreenConfigurationByAssemblyId(lC); err != nil {
					showErrorMessageBox(nil, err_msg.ConfigurationSaveErr)
					return
				}
				scaledPixmap := newScaledPixmap.Scaled2(
					rootWidth-210,
					rootHeight,
					core.Qt__IgnoreAspectRatio,
					core.Qt__SmoothTransformation,
				)
				backgroundLabel.SetScaledContents(true)
				backgroundLabel.SetPixmap(scaledPixmap)
				palette := mainWindow.Palette()
				brush := gui.NewQBrush3(gui.NewQColor3(0, 0, 0, 0), core.Qt__TexturePattern)
				brush.SetTexture(scaledPixmap)
				palette.SetBrush(gui.QPalette__Window, brush)
				mainWindow.SetPalette(palette)
				mainWindow.Repaint()
				showSuccessMessageBox(nil, "保存成功")
			})

			fileDialog.Exec()
		})

		//背景图纯色
		addBackColorAction.ConnectTriggered(func(bool) {
			dialog := widgets.NewQDialog(nil, 0)
			backColorLayout := widgets.NewQGridLayout(nil)
			// 创建按钮并将其添加到布局中
			backColorButton := widgets.NewQPushButton2("选择颜色", nil)
			colorStr := ""
			backColorButton.ConnectClicked(func(bool) {
				color := widgets.QColorDialog_GetColor(gui.NewQColor3(0, 0, 0, 255), nil, "选择颜色", 0)
				if color.IsValid() {
					//centralWidget.SetStyleSheet("background-color:" + color.Name())
					pixmap := gui.NewQPixmap()
					pixmap.Fill(color)
					// 设置像素图作为新的背景
					backgroundLabel.SetPixmap(pixmap)
					palette := mainWindow.Palette()
					brush := gui.NewQBrush3(color, core.Qt__SolidPattern)
					palette.SetBrush(gui.QPalette__Window, brush)
					mainWindow.SetPalette(palette)
					mainWindow.Repaint()
				}
				colorStr = color.Name()
			})
			backColorExplainLabel := widgets.NewQLabel(nil, 0)
			backColorExplainLabel.SetText("颜色设置后会清空背景图！")
			backColorLayout.AddWidget3(backColorButton, 0, 0, 1, 2, 0)
			backColorLayout.AddWidget3(backColorExplainLabel, 1, 0, 1, 2, 0)
			okButton := widgets.NewQPushButton2("确定", dialog)
			backColorLayout.AddWidget3(okButton, 2, 0, 1, 1, 0)
			okButton.ConnectClicked(func(bool) {
				var lc database.LargeScreenConfiguration
				lc.MeetingId = executeMeetingId
				lc.TicketId = ticketId
				lc.TypeNum = 0
				aId := backgroundLabel.ObjectName()
				if aId == "" {
					aId = xid.New().String()
					backgroundLabel.SetObjectName(aId)
				}
				lc.AssemblyId = aId
				lc.BackColor = colorStr
				lc.LargeScreenId = displayLargeScreenType
				if err := database.CreateOrUpdateLargeScreenConfigurationByAssemblyId(lc); err != nil {
					showErrorMessageBox(nil, err_msg.ConfigurationSaveErr)
					return
				}
				showSuccessMessageBox(nil, "修改成功")
				dialog.Reject()
			})
			cancelButton := widgets.NewQPushButton2("取消", dialog)
			backColorLayout.AddWidget3(cancelButton, 2, 1, 1, 1, 0)
			cancelButton.ConnectClicked(func(bool) {
				backGroundImage := database.FindAllAssembliesByMeetingIdAndTypeNum(executeMeetingId, 0, displayLargeScreenType)
				if len(backGroundImage) == 0 {
					backGroundImage = append(backGroundImage, database.FindDefaultBackGround(displayLargeScreenType))
					os.WriteFile(backGroundImage[0].BackPictureDir, backGroundImage[0].BackPictureData, 0644)
					backgroundPixmap := gui.NewQPixmap3(backGroundImage[0].BackPictureDir, "", core.Qt__AutoColor)
					scaledPixmap := backgroundPixmap.Scaled2(
						rootWidth-210,
						rootHeight,
						core.Qt__IgnoreAspectRatio,
						core.Qt__SmoothTransformation,
					)
					backgroundLabel.SetScaledContents(true)
					backgroundLabel.SetPixmap(scaledPixmap)
					// 设置QPalette对象
					palette := mainWindow.Palette()
					brush := gui.NewQBrush3(gui.NewQColor3(0, 0, 0, 0), core.Qt__TexturePattern)
					brush.SetTexture(scaledPixmap)
					palette.SetBrush(gui.QPalette__Window, brush)
					mainWindow.SetPalette(palette)
				} else {
					backgroundLabel.SetObjectName(backGroundImage[0].AssemblyId)
					if backGroundImage[0].BackPictureDir != "" {
						os.WriteFile(backGroundImage[0].BackPictureDir, backGroundImage[0].BackPictureData, 0644)
						backgroundPixmap := gui.NewQPixmap3(backGroundImage[0].BackPictureDir, "", core.Qt__AutoColor)
						scaledPixmap := backgroundPixmap.Scaled2(
							rootWidth-210,
							rootHeight,
							core.Qt__IgnoreAspectRatio,
							core.Qt__SmoothTransformation,
						)
						backgroundLabel.SetScaledContents(true)
						backgroundLabel.SetPixmap(scaledPixmap)
						// 设置QPalette对象
						palette := mainWindow.Palette()
						brush := gui.NewQBrush3(gui.NewQColor3(0, 0, 0, 0), core.Qt__TexturePattern)
						brush.SetTexture(scaledPixmap)
						palette.SetBrush(gui.QPalette__Window, brush)
						mainWindow.SetPalette(palette)
					}
					if backGroundImage[0].BackColor != "" {
						style := fmt.Sprintf("background-color:%s;", backGroundImage[0].BackColor)
						backgroundLabel.SetStyleSheet(style)
					}
				}
				dialog.Reject()
			})
			dialog.SetLayout(backColorLayout)
			dialog.SetMinimumSize2(300, 200)
			dialog.SetWindowTitle("选择背景纯色")
			// 显示对话框
			dialog.Show()
		})

		//退出关闭
		returnBackAction.ConnectTriggered(func(bool) {
			fmt.Println("退出关闭")
			window.Close()
		})

		mainWindow.Update()

	})

	//mainWindow.SetStyleSheet(`border: 1px solid #8f8f91;`)
	//mainWindow.Show()

	windowLayout.AddWidget(mainWindow, 0, 0)
	windowCentralWidget.SetLayout(windowLayout)
	window.SetCentralWidget(windowCentralWidget)
	// 显示标签
	window.Show()
	app.Exec()

}

func ImageLabelAction(imgLabel *widgets.QLabel, pixmap *gui.QPixmap) {
	//双击修改
	imgLabel.ConnectMouseDoubleClickEvent(func(event *gui.QMouseEvent) {
		dialog := widgets.NewQDialog(nil, 0)
		imgOperateLayout := widgets.NewQGridLayout(nil)
		imgOperateLayout.SetColumnStretch(1, 2)
		var imgFiles string
		imgFileLabel := widgets.NewQLabel2("选择的图片文件：", nil, 0)
		uploadImgBtn := widgets.NewQPushButton2("上传图片文件", nil)
		uploadImgBtn.SetStyleSheet("QPushButton { background-color: #5D9CEC; color: white; padding: 5px; border-radius: 5px; } QPushButton:hover { background-color: #4A89DC; }")
		imgOperateLayout.AddWidget3(imgFileLabel, 0, 0, 1, 1, 0)
		uploadImgBtn.ConnectClicked(func(checked bool) {
			fileDialog := widgets.NewQFileDialog(nil, 0)
			fileDialog.SetAcceptMode(widgets.QFileDialog__AcceptOpen)
			fileDialog.SetFileMode(widgets.QFileDialog__ExistingFile)
			fileDialog.SetNameFilter("图片文件 (*.jpg *.jpeg *.png *.bmp *.gif)")
			fileDialog.ConnectFileSelected(func(file string) {
				// 获取所选文件的大小
				fileInfo := core.NewQFileInfo3(file)
				fileSize := int(fileInfo.Size())
				// 限制图片大小为 10 MB
				maxFileSize := int64(10 * 1024 * 1024)
				if int64(fileSize) > maxFileSize {
					showErrorMessageBox(nil, err_msg.PictureTooBigErr)
					return
				}
				imgFiles = file
				imgFileLabel.SetText("选择的图片文件：" + imgFiles)
				// 缩放图片
				newScaledPixmap := gui.NewQPixmap3(imgFiles, "", core.Qt__AutoColor)
				scaleFactor := float64(imgLabel.Width()) / float64(newScaledPixmap.Width())
				if float64(newScaledPixmap.Height())*scaleFactor > float64(imgLabel.Height()) {
					scaleFactor = float64(imgLabel.Height()) / float64(newScaledPixmap.Height())
				}
				scaledWidth := int(float64(newScaledPixmap.Width()) * scaleFactor)
				scaledHeight := int(float64(newScaledPixmap.Height()) * scaleFactor)
				scaledPixmap := newScaledPixmap.Scaled2(
					scaledWidth,
					scaledHeight,
					core.Qt__IgnoreAspectRatio,
					core.Qt__SmoothTransformation,
				)

				// 更新界面
				imgLabel.SetPixmap(scaledPixmap)
				palette := imgLabel.Palette()
				brush := gui.NewQBrush3(gui.NewQColor3(0, 0, 0, 0), core.Qt__TexturePattern)
				brush.SetTexture(scaledPixmap)
				palette.SetBrush(gui.QPalette__Window, brush)
				imgLabel.SetPalette(palette)
				imgLabel.Repaint()
			})
			fileDialog.Exec()
		})
		imgOperateLayout.AddWidget3(uploadImgBtn, 0, 1, 1, 2, 0)

		okButton := widgets.NewQPushButton2("确定", dialog)
		imgOperateLayout.AddWidget3(okButton, 2, 0, 1, 1, 0)

		cancelButton := widgets.NewQPushButton2("取消", dialog)
		imgOperateLayout.AddWidget3(cancelButton, 2, 2, 1, 1, 0)

		dialog.SetLayout(imgOperateLayout)
		dialog.SetMinimumSize2(300, 200)
		dialog.SetWindowTitle("修改图片及大小")
		okButton.ConnectClicked(func(bool) {
			var lC database.LargeScreenConfiguration
			if imgFiles != "" {
				file, imgErr := ioutil.ReadFile(imgFiles)
				if imgErr != nil {
					showErrorMessageBox(nil, err_msg.ImageReadErr)
					return
				}
				lC.PictureDir = imgFiles
				lC.PictureData = file
			} else {
				showErrorMessageBox(nil, err_msg.OperateLargeScreenImageErr)
				return
			}
			lC.AssemblyId = imgLabel.ObjectName()
			lC.TypeNum = 1
			lC.LargeScreenId = displayLargeScreenType

			lC.PictureWidth = pixmap.Width()
			lC.PictureHeight = pixmap.Height()
			if err := database.CreateOrUpdateLargeScreenConfigurationByAssemblyId(lC); err != nil {
				showErrorMessageBox(nil, err_msg.ConfigurationSaveErr)
				return
			}
			showSuccessMessageBox(nil, "保存成功")
			imgLabel.Show()
			dialog.Accept()
		})
		cancelButton.ConnectClicked(func(bool) {
			// 关闭对话框，并将 textLabel 显示出来
			img := database.FindByAssemblyId(imgLabel.ObjectName())
			// 缩放图片
			newScaledPixmap := gui.NewQPixmap3(img.PictureDir, "", core.Qt__AutoColor)
			scaleFactor := float64(imgLabel.Width()) / float64(newScaledPixmap.Width())
			if float64(newScaledPixmap.Height())*scaleFactor > float64(imgLabel.Height()) {
				scaleFactor = float64(imgLabel.Height()) / float64(newScaledPixmap.Height())
			}
			scaledWidth := int(float64(newScaledPixmap.Width()) * scaleFactor)
			scaledHeight := int(float64(newScaledPixmap.Height()) * scaleFactor)
			scaledPixmap := newScaledPixmap.Scaled2(
				scaledWidth,
				scaledHeight,
				core.Qt__IgnoreAspectRatio,
				core.Qt__SmoothTransformation,
			)

			// 更新界面
			imgLabel.SetPixmap(scaledPixmap)
			palette := imgLabel.Palette()
			brush := gui.NewQBrush3(gui.NewQColor3(0, 0, 0, 0), core.Qt__TexturePattern)
			brush.SetTexture(scaledPixmap)
			palette.SetBrush(gui.QPalette__Window, brush)
			imgLabel.SetPalette(palette)
			imgLabel.Repaint()
			dialog.Reject()
			imgLabel.Show()
		})
		// 显示对话框
		dialog.Show()
	})
	//拖拽
	//设置标签为可拖拽
	imgLabel.SetAttribute(core.Qt__WA_MouseTracking, true)
	imgLabel.SetAttribute(core.Qt__WA_TransparentForMouseEvents, false)
	imgLabel.SetMouseTracking(true)
	// 标签的位置偏移量
	offsetX, offsetY := 0, 0
	// 重写标签的 mousePressEvent() 函数，以便在鼠标点击标签时设置标签为可拖拽
	imgLabel.ConnectMousePressEvent(func(event *gui.QMouseEvent) {
		if event.Button() == core.Qt__LeftButton {
			offsetX = event.Pos().X() + 210
			offsetY = event.Pos().Y()
		}
	})
	// 重写标签的 mouseMoveEvent() 函数，以便在拖拽标签时移动标签的位置
	imgLabel.ConnectMouseMoveEvent(func(event *gui.QMouseEvent) {
		if event.Buttons()&core.Qt__LeftButton > 0 {
			x := event.GlobalPos().X() - offsetX
			y := event.GlobalPos().Y() - offsetY
			imgLabel.Move2(x, y)
			var lc database.LargeScreenConfiguration
			//lc.MeetingId = executeMeetingId
			lc.AssemblyId = imgLabel.ObjectName()
			lc.ScreenRowX = float64(x) / float64(rootWidth)
			lc.ScreenColY = float64(y) / float64(rootHeight)
			lc.TypeNum = 1
			lc.LargeScreenId = displayLargeScreenType
			err := database.CreateOrUpdateLargeScreenConfigurationByAssemblyId(lc)
			if err != nil {
				showErrorMessageBox(nil, err_msg.OperateLargeScreenTextErr)
			}
		}
	})

	// 获取当前的图片尺寸
	currentWidth := float64(pixmap.Width())
	currentHeight := float64(pixmap.Height())

	// 添加鼠标滚轮事件侦听器
	imgLabel.ConnectWheelEvent(func(event *gui.QWheelEvent) {
		// 获取鼠标滚轮事件中的滚动角度
		angle := event.AngleDelta().Y() / 120.0
		// 计算缩放比例
		scaleFactor := 1.0
		if angle > 0 {
			scaleFactor = scaleFactor + 0.1
		} else if angle < 0 {
			scaleFactor = scaleFactor - 0.1
		}
		// 计算新的图片尺寸
		newWidth := currentWidth * scaleFactor
		newHeight := currentHeight * scaleFactor
		// 获取当前的图片尺寸
		currentWidth = newWidth
		currentHeight = newHeight
		// 缩放图片
		newPixmap := pixmap.Scaled2(int(newWidth), int(newHeight), core.Qt__KeepAspectRatio, core.Qt__SmoothTransformation)
		imgLabel.SetPixmap(newPixmap)
		var lc database.LargeScreenConfiguration
		//lc.MeetingId = executeMeetingId
		lc.AssemblyId = imgLabel.ObjectName()
		lc.TypeNum = 1
		lc.PictureWidth = int(newWidth)
		lc.PictureHeight = int(newHeight)
		lc.LargeScreenId = displayLargeScreenType
		err := database.CreateOrUpdateLargeScreenConfigurationByAssemblyId(lc)
		if err != nil {
			showErrorMessageBox(nil, err_msg.OperateLargeScreenTextErr)
		}
	})
}

//文本事件封装
func TextLabelAction(textLabel *widgets.QLabel, centralWidget *widgets.QWidget) {
	// 双击编辑文本
	textLabel.ConnectMouseDoubleClickEvent(func(event *gui.QMouseEvent) {
		lineEdit := widgets.NewQLineEdit(centralWidget)
		lineEdit.SetText(textLabel.Text())
		lineEdit.Move2(event.GlobalPos().X(), event.GlobalPos().Y())
		lineEdit.SetFixedSize2(textLabel.Width(), textLabel.Height())
		lineEdit.SetStyleSheet("border: 1px solid black;")
		lineEdit.Show()
		textLabel.Hide()

		lineEdit.ConnectEditingFinished(func() {
			text := lineEdit.Text()
			var lc database.LargeScreenConfiguration
			//lc.MeetingId = executeMeetingId
			lc.AssemblyId = textLabel.ObjectName()
			lc.TextContent = text
			//lc.ScreenRowX = textLabel.X()
			//lc.ScreenRowY = textLabel.Y()
			lc.TypeNum = 2
			lc.LargeScreenId = displayLargeScreenType
			qFont := textLabel.Font()
			//lc.TextSizeInt = qFont.PointSize()
			//lc.TextStyle = qFont.Family()
			fontMetrics := gui.NewQFontMetrics(qFont)
			textWidth := fontMetrics.HorizontalAdvance(textLabel.Text(), len(textLabel.Text()))
			textHeight := fontMetrics.Height() + 2
			textLabel.SetFixedSize2(textWidth, textHeight)
			lc.TextWidth = textWidth
			lc.TextHeight = textHeight
			var isWeight bool
			if qFont.Weight() == 75 {
				isWeight = true
			} else {
				isWeight = false
			}
			lc.IsTextWeight = isWeight
			lc.IsTextItalic = qFont.Italic()
			lc.IsTextUnder = qFont.Underline()
			lc.IsTextStrikeOut = qFont.StrikeOut()
			err := database.CreateOrUpdateLargeScreenConfigurationByAssemblyId(lc)
			if err != nil {
				showErrorMessageBox(nil, err_msg.OperateLargeScreenTextErr)
			}
			textLabel.SetText(text)
			textLabel.Show()
			lineEdit.Close()
		})

		//双击修改字体
		lineEdit.ConnectMouseDoubleClickEvent(func(event *gui.QMouseEvent) {
			textAssembly := database.FindByAssemblyId(textLabel.ObjectName())
			dialog := widgets.NewQDialog(nil, 0)
			textSizeLayout := widgets.NewQGridLayout(nil)
			textSizeLayout.SetColumnStretch(3, 1)
			textSizeTitle := widgets.NewQLabel(nil, 0)
			textSizeTitle.SetText("修改字体大小：")
			textSizeLineEdit := widgets.NewQLineEdit(nil)
			textDefaultSize := textLabel.Font().PointSize()
			if reflect.ValueOf(&textAssembly).IsNil() {
				textDefaultSize = textAssembly.TextSizeInt
			}
			textSizeLineEdit.SetText(fmt.Sprintf("%d", textDefaultSize))
			textSizeLayout.AddWidget3(textSizeTitle, 0, 0, 1, 1, 0)
			textSizeLayout.AddWidget3(textSizeLineEdit, 0, 1, 1, 1, 0)

			textStyleTitle := widgets.NewQLabel(nil, 0)
			textStyleTitle.SetText("选择字体：")
			styleBox := widgets.NewQComboBox(nil)
			styleBox.AddItems([]string{
				"宋体",
				"微软雅黑",
				"黑体",
				"隶书",
				"等线",
				"楷体",
			})
			styleBox.SetCurrentText(textAssembly.TextStyle)
			textSizeLayout.AddWidget3(textStyleTitle, 1, 0, 1, 1, 0)
			textSizeLayout.AddWidget3(styleBox, 1, 1, 1, 1, 0)

			isWeightBox := widgets.NewQCheckBox(nil)
			isWeightBox.SetText("是否加粗")
			isWeightBox.SetChecked(textAssembly.IsTextWeight)
			isWeightBox.ConnectClicked(func(bool) {
				font := textLabel.Font()
				if isWeightBox.IsChecked() {
					font.SetWeight(75)
				} else {
					font.SetWeight(50)
				}
				textLabel.SetFont(font)
				textLabel.Update()
			})
			textSizeLayout.AddWidget3(isWeightBox, 2, 0, 1, 2, 0)

			isItalicBox := widgets.NewQCheckBox(nil)
			isItalicBox.SetText("是否斜体")
			isItalicBox.SetChecked(textAssembly.IsTextItalic)
			textSizeLayout.AddWidget3(isItalicBox, 3, 0, 1, 2, 0)

			isUnderlineBox := widgets.NewQCheckBox(nil)
			isUnderlineBox.SetText("是否带下划线")
			isUnderlineBox.SetChecked(textAssembly.IsTextUnder)
			textSizeLayout.AddWidget3(isUnderlineBox, 4, 0, 1, 2, 0)

			isStrikeOutBox := widgets.NewQCheckBox(nil)
			isStrikeOutBox.SetText("是否带删除线")
			isStrikeOutBox.SetChecked(textAssembly.IsTextStrikeOut)
			textSizeLayout.AddWidget3(isStrikeOutBox, 5, 0, 1, 2, 0)

			// 创建滑块控件
			opacitySlider := widgets.NewQSlider2(core.Qt__Horizontal, nil)
			opacitySlider.SetMinimum(0)
			opacitySlider.SetMaximum(100)
			opacitySlider.SetTickInterval(10)
			opacitySlider.SetValue(50)

			opacityLabel := widgets.NewQLabel(nil, core.Qt__Widget)
			opacityLabel.SetStyleSheet("background-color: white; opacity: 0.5;")

			colorStr := ""
			// 创建按钮并将其添加到布局中
			ChooseColorButton := widgets.NewQPushButton2("选择颜色", nil)
			ChooseColorButton.ConnectClicked(func(checked bool) {
				color := widgets.QColorDialog_GetColor(gui.NewQColor3(0, 0, 0, 255), nil, "选择颜色", 0)
				if color.IsValid() {
					colorStr = color.Name()
					// 更新标签的颜色和透明度
					textColor := gui.NewQColor6(colorStr)
					textColor.SetAlphaF(float64(opacitySlider.Value()) / 100.0)
					textLabel.SetStyleSheet(fmt.Sprintf("color: rgba(%d, %d, %d, %f);", textColor.Red(), textColor.Green(), textColor.Blue(), textColor.AlphaF()))
					opacityLabelStyle := fmt.Sprintf("background-color: rgba(%d, %d, %d, %f); color: rgba(%d, %d, %d, %f);", textColor.Red(), textColor.Green(), textColor.Blue(), textColor.AlphaF(), textColor.Red(), textColor.Green(), textColor.Blue(), textColor.AlphaF())
					opacityLabel.SetStyleSheet(opacityLabelStyle)
					opacityLabel.Repaint()
				}
			})
			textSizeLayout.AddWidget3(ChooseColorButton, 4, 2, 1, 1, 0)

			var textAlphaF float64
			opacitySlider.ConnectValueChanged(func(value int) {
				opacity := float64(value) / 100.0
				textColor := gui.NewQColor6(colorStr)
				textColor.SetAlphaF(opacity)
				textAlphaF = textColor.AlphaF()
				textLabel.SetStyleSheet(fmt.Sprintf("color: rgba(%d, %d, %d, %f);", textColor.Red(), textColor.Green(), textColor.Blue(), textColor.AlphaF()))
				opacityLabelStyle := fmt.Sprintf("background-color: rgba(%d, %d, %d, %f); color: rgba(%d, %d, %d, %f);", textColor.Red(), textColor.Green(), textColor.Blue(), textColor.AlphaF(), textColor.Red(), textColor.Green(), textColor.Blue(), textColor.AlphaF())
				opacityLabel.SetStyleSheet(opacityLabelStyle)
				opacityLabel.Repaint()
			})

			textSizeLayout.AddWidget3(opacityLabel, 6, 0, 1, 1, 0)
			textSizeLayout.AddWidget3(opacitySlider, 6, 1, 1, 2, 0)

			okButton := widgets.NewQPushButton2("确定", dialog)
			textSizeLayout.AddWidget3(okButton, 7, 0, 1, 1, 0)

			cancelButton := widgets.NewQPushButton2("取消", dialog)
			textSizeLayout.AddWidget3(cancelButton, 7, 2, 1, 1, 0)

			dialog.SetLayout(textSizeLayout)
			dialog.SetMinimumSize2(300, 200)
			dialog.SetWindowTitle("修改字体及大小")
			okButton.ConnectClicked(func(checked bool) {
				var lc database.LargeScreenConfiguration
				//lc.MeetingId = executeMeetingId
				lc.AssemblyId = textLabel.ObjectName()
				//lc.ScreenRowX = textLabel.X()
				//lc.ScreenRowY = textLabel.Y()
				lc.TypeNum = 2
				//lc.TextContent = textLabel.Text()
				//lc.TextWidth = textLabel.Width()
				//lc.TextHeight = textLabel.Height()
				// 获取用户选择的字体大小和字体样式
				fontPointSizeStr := textSizeLineEdit.Text()
				re := regexp.MustCompile(`^[1-9]\d?$|^100$`)
				if !re.MatchString(fontPointSizeStr) {
					widgets.QMessageBox_Information(nil, "提示", "请输入有效的字体大小(1~100)", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
					return
				}
				fontPointSize, _ := strconv.Atoi(fontPointSizeStr)
				fontFamily := styleBox.CurrentText()

				// 创建一个 QFont 对象
				font := textLabel.Font()
				font.SetPointSize(fontPointSize)
				lc.TextSizeInt = fontPointSize
				font.SetFamily(fontFamily)
				lc.TextStyle = fontFamily
				lc.TextColor = colorStr
				fontMetrics := gui.NewQFontMetrics(font)
				textWidth := fontMetrics.HorizontalAdvance(textLabel.Text(), len(textLabel.Text()))
				textHeight := fontMetrics.Height() + 2
				textLabel.SetFixedSize2(textWidth, textHeight)
				lc.TextWidth = textWidth
				lc.TextHeight = textHeight
				lc.LargeScreenId = displayLargeScreenType
				lc.Opacity = textAlphaF
				// 根据用户选择设置字体样式
				if isWeightBox.IsChecked() {
					font.SetWeight(75)
				} else {
					font.SetWeight(50)
				}
				lc.IsTextWeight = isWeightBox.IsChecked()
				font.SetItalic(isItalicBox.IsChecked())
				lc.IsTextItalic = isItalicBox.IsChecked()
				font.SetUnderline(isUnderlineBox.IsChecked())
				lc.IsTextUnder = isUnderlineBox.IsChecked()
				font.SetStrikeOut(isStrikeOutBox.IsChecked())
				lc.IsTextStrikeOut = isStrikeOutBox.IsChecked()

				// 将新的字体应用到标签
				textLabel.SetFont(font)

				err := database.CreateOrUpdateLargeScreenConfigurationByAssemblyId(lc)
				if err != nil {
					showErrorMessageBox(nil, err_msg.OperateLargeScreenTextErr)
				}
				// 关闭对话框
				dialog.Accept()
			})
			cancelButton.ConnectClicked(func(checked bool) {
				// 关闭对话框，并将 textLabel 显示出来
				text := database.FindByAssemblyId(textLabel.ObjectName())
				textLabel.SetText(text.TextContent)
				textLabel.SetMinimumSize2(200, 50)
				textLabel.SetFixedSize2(text.TextWidth, text.TextHeight)
				textLabel.SetObjectName(text.AssemblyId)
				layout := centralWidget.Layout()
				layout.AddWidget(textLabel)
				x, _ := strconv.Atoi(fmt.Sprintf("%1.0f", text.ScreenRowX*float64(rootWidth)))
				y, _ := strconv.Atoi(fmt.Sprintf("%1.0f", text.ScreenColY*float64(rootHeight)))
				textLabel.Move2(x, y)
				textColor := gui.NewQColor6(text.TextColor)
				textColor.SetAlphaF(text.Opacity)
				style := fmt.Sprintf("color: rgba(%d, %d, %d, %f);", textColor.Red(), textColor.Green(), textColor.Blue(), textColor.AlphaF())
				textLabel.SetStyleSheet(style)
				// 创建一个 QFont 对象
				font := textLabel.Font()
				font.SetPointSize(text.TextSizeInt)
				font.SetFamily(text.TextStyle)
				if text.IsTextWeight {
					font.SetWeight(75)
				} else {
					font.SetWeight(50)
				}
				font.SetItalic(text.IsTextItalic)
				font.SetUnderline(text.IsTextUnder)
				font.SetStrikeOut(text.IsTextStrikeOut)
				// 将新的字体应用到标签
				textLabel.SetFont(font)
				dialog.Reject()
				textLabel.Show()
			})
			// 显示对话框
			dialog.Show()

			// 隐藏 textLabel
			textLabel.Hide()
		})
	})

	//拖拽
	//设置标签为可拖拽
	textLabel.SetAttribute(core.Qt__WA_MouseTracking, true)
	textLabel.SetAttribute(core.Qt__WA_TransparentForMouseEvents, false)
	textLabel.SetMouseTracking(true)
	// 标签的位置偏移量
	offsetX, offsetY := 0, 0
	// 重写标签的 mousePressEvent() 函数，以便在鼠标点击标签时设置标签为可拖拽
	textLabel.ConnectMousePressEvent(func(event *gui.QMouseEvent) {
		if event.Button() == core.Qt__LeftButton {
			offsetX = event.Pos().X() - 210
			offsetY = event.Pos().Y()
		}
	})
	// 重写标签的 mouseMoveEvent() 函数，以便在拖拽标签时移动标签的位置
	textLabel.ConnectMouseMoveEvent(func(event *gui.QMouseEvent) {
		if event.Buttons()&core.Qt__LeftButton > 0 {
			x := event.GlobalPos().X() - offsetX
			y := event.GlobalPos().Y() - offsetY
			textLabel.Move2(x, y)
			var lc database.LargeScreenConfiguration
			//lc.MeetingId = executeMeetingId
			lc.AssemblyId = textLabel.ObjectName()
			//lc.TextContent = textLabel.Text()
			lc.ScreenRowX = float64(x)
			lc.ScreenColY = float64(y)
			lc.TypeNum = 2
			lc.LargeScreenId = displayLargeScreenType
			qFont := textLabel.Font()
			//lc.TextSizeInt = qFont.PointSize()
			//lc.TextStyle = qFont.Family()
			//lc.TextWidth = textLabel.Width()
			//lc.TextHeight = textLabel.Height()
			var isWeight bool
			if qFont.Weight() == 75 {
				isWeight = true
			} else {
				isWeight = false
			}
			lc.IsTextWeight = isWeight
			lc.IsTextItalic = qFont.Italic()
			lc.IsTextUnder = qFont.Underline()
			lc.IsTextStrikeOut = qFont.StrikeOut()
			err := database.CreateOrUpdateLargeScreenConfigurationByAssemblyId(lc)
			if err != nil {
				showErrorMessageBox(nil, err_msg.OperateLargeScreenTextErr)
			}
		}
	})

}

func showErrorMessageBox(parent *widgets.QWidget, text string) {
	errorBox := widgets.NewQMessageBox(parent)
	errorBox.SetIcon(widgets.QMessageBox__Warning)
	errorBox.SetWindowTitle("错误")
	errorBox.SetText(text)
	errorBox.SetStandardButtons(widgets.QMessageBox__Ok)

	// 设置样式
	errorBox.SetStyleSheet(`
        QMessageBox {
            background-color: #F8F8F8;
        }
        QMessageBox QLabel {
            color: #333333;
            font-size: 14px;
        }
        QMessageBox QPushButton {
            background-color: #EBEBEB;
            border: 1px solid #CCCCCC;
            border-radius: 4px;
            color: #333333;
            font-size: 14px;
            padding: 5px 10px;
            margin-right: 10px;
        }
        QMessageBox QPushButton:hover {
            background-color: #E2E2E2;
        }
    `)

	errorBox.Exec()
}

func showSuccessMessageBox(parent *widgets.QWidget, text string) {
	successBox := widgets.NewQMessageBox(parent)
	successBox.SetIcon(widgets.QMessageBox__Information)
	successBox.SetWindowTitle("成功")
	successBox.SetText(text)
	successBox.SetStandardButtons(widgets.QMessageBox__Ok)
	successBox.SetStyleSheet(`
	QMessageBox {
		background-color: #EAF7D8;
		color: #2F9D27;
	}
	QMessageBox QPushButton {
		background-color: #FFFFFF;
		border: 1px solid #2F9D27;
		color: #2F9D27;
	}
	QMessageBox QPushButton:hover {
		background-color: #EAF7D8;
		color: #2F9D27;
	}
`)
	successBox.Exec()

}

func TestListWidgetChoose(t *testing.T) {
	// 创建Qt应用程序
	app := widgets.NewQApplication(len(os.Args), os.Args)
	// 创建主窗口
	window := widgets.NewQMainWindow(nil, 0)
	window.SetFixedSize2(800, 600)
	windowCentralWidget := widgets.NewQWidget(nil, 0)
	window.SetWindowTitle("大屏展示配置")
	windowLayout := widgets.NewQHBoxLayout()

	ticketImageList := widgets.NewQListWidget(nil)
	ticketImageList.SetViewMode(widgets.QListView__IconMode)
	ticketImageList.SetIconSize(core.NewQSize2(200, 200))
	ticketImageList.SetFixedWidth(210)
	ticketImageList.SetResizeMode(widgets.QListView__Adjust)
	ticketPics, _ := database.FindTicketsByMeetingId(executeMeetingId)

	for _, thisT := range ticketPics {
		os.WriteFile(thisT.PictureDir, thisT.PictureData, 0644)
		pixmap := gui.NewQPixmap3(thisT.PictureDir, "", 0)
		item := widgets.NewQListWidgetItem2("", ticketImageList, 0)
		item.SetIcon(gui.NewQIcon2(pixmap))
		font := item.Font()
		font.SetPointSize(9)
		item.SetFont(font)
		fm := gui.NewQFontMetrics(font)
		elidedText := fm.ElidedText(thisT.ConferenceTitleName, core.Qt__ElideRight, 180, 0)
		item.SetText(elidedText)
		item.SetTextAlignment(int(core.Qt__AlignCenter))
		item.SetData(int(core.Qt__UserRole), core.NewQVariant1(thisT.TicketId))
	}
	windowLayout.AddWidget(ticketImageList, 0, 0)
	mainWindow := widgets.NewQWidget(nil, 0)
	ticketImageList.ConnectItemClicked(func(item *widgets.QListWidgetItem) {
		ticketId := item.Data(int(core.Qt__UserRole)).ToString()
		l := widgets.NewQLabel(nil, 0)
		l.SetText(ticketId)

		layout := widgets.NewQVBoxLayout()
		layout.AddWidget(l, 0, 0)
		mainWindow.SetLayout(layout)
		mainWindow.Repaint()
	})
	//mainWindow.SetStyleSheet(`border: 1px solid #8f8f91;`)
	windowLayout.AddWidget(mainWindow, 0, 0)
	windowCentralWidget.SetLayout(windowLayout)
	window.SetCentralWidget(windowCentralWidget)
	// 显示标签
	window.Show()
	app.Exec()
}

func TestSetItemType(t *testing.T) {
	// 创建Qt应用程序
	app := widgets.NewQApplication(len(os.Args), os.Args)
	//ticketId := "cgn45s3rm5p0j51tt0eg"
	// 创建主窗口
	window := widgets.NewQMainWindow(nil, 0)
	window.SetFixedSize2(800, 600)
	windowCentralWidget := widgets.NewQWidget(nil, 0)
	window.SetWindowTitle("表决事项规则配置")
	windowLayout := widgets.NewQVBoxLayout()

	//ts, err := database.FindTicketItemByTicketId(ticketId)
	//if err != nil {
	//	showErrorMessageBox(nil, err_msg.AnalysisItemTitleErr)
	//}

	titleLabel := widgets.NewQLabel(nil, 0)
	titleLabel.SetText("表决事项")

	//for _, item := range ts {
	//
	//}
	windowCentralWidget.SetLayout(windowLayout)
	window.SetCentralWidget(windowCentralWidget)
	// 显示标签
	window.Show()
	app.Exec()
}
