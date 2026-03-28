#ifndef _HPRT_ERROR_H_
#define _HPRT_ERROR_H_

#include "hprt_define.h"

/*///////////////////////////////////////////////////////////////////
// "error_code" definitions
///////////////////////////////////////////////////////////////////*/
#define HPRT_SUCCESS						0		// no error

/*///////////////////////////////////////////////////////////////////
// common  error code definitions
///////////////////////////////////////////////////////////////////*/

#define HPRT_E_INVALID_PARAMETER			-1
#define HPRT_E_NOT_ENOUGH_BUFFER			-2
#define HPRT_E_INVALID_MODEL_TYPE			-3
#define HPRT_E_NOT_SUPPORT					-4
//#define HPRT_E_NOT_OPEN						-5
#define HPRT_E_BAD_HANDLE					-6
#define HPRT_E_NOT_IMPLEMENTED				-7
#define HPRT_E_INVALID_MODEL				-8
#define HPRT_E_NOT_ENOUGH_MEMORY			-9
#define HPRT_E_NOT_TARGET_PRINTER           -10
#define HPRT_E_INVALID_ENVIRONMENT          -11     /* most used in Linux */
//add cah 2024.3.13 
#define HPRT_E_LOAD_CONFIG_FAILED			 -12 
//add cah 2024.4.11
#define HPRT_E_PRINTERTYPE_ERROR			-13

/* Image Error */
#define HPRT_E_OPEN_FILE_ERROR				-20
#define HPRT_E_LOAD_IMAGE_ERROR				-21
#define HPRT_E_ANALYSIS_IMAGE_ERROR			-22
#define HPRT_E_IMAGE_BAD_SIZE				-25

#define HPRT_E_INVALID_DATA                 -30     /* Data is not correct */
#define HPRT_E_OPEN_LOG_ERROR				-31

#define HPRT_E_BASE                         -100
/* Driver */
#define HPRT_E_DRIVER_INCORRECT_DATA        -101
#define HPRT_E_DRIVER_PRINTER_STATE_ERROR   -102
#define HPRT_E_DRIVER_TCP_NOFOUND           -103
/* IO */

/* io setting error */
#define HPRT_E_IO_ERROR						-300
#define HPRT_E_IO_INVALID_SETTING			-301    
#define HPRT_E_IO_NAME_TOO_LONG				-302
#define HPRT_E_IO_OS_VERSION_TOO_LOW		-304

#define HPRT_E_IO_INVALID_HANDLE			 -308
#define HPRT_E_IO_PORT_NOT_OPEN				 -309
#define HPRT_E_IO_PORT_ALREADY_OPEN			 -310

/* io open error */
#define HPRT_E_IO_OPEN_FAILED				 -311   
/* io attr get/set error */
#define HPRT_E_IO_GETATTR_ERROR				 -312    
#define HPRT_E_IO_SETATTR_ERROR				 -313
/* io write error */
#define HPRT_E_IO_WRITE_FAILED				 -321    
#define HPRT_E_IO_WRITE_TIMEOUT				 -322    
/* io readerror */
#define HPRT_E_IO_READ_FAILED				 -331    
#define HPRT_E_IO_READ_TIMEOUT				 -332  
/* io flush error */
#define HPRT_E_IO_FLUSH_FAILED				 -341
/* serial port error */
#define HPRT_E_IO_SERIAL_INVALID_BAUDRATE	 -351
#define HPRT_E_IO_SERIAL_INVALID_HANDSHAKE	 -352
#define HPRT_E_IO_SERIAL_INVALID_PARITY		 -353
#define HPRT_E_IO_SERIAL_INVALID_BYTESIZE	 -354
#define HPRT_E_IO_SERIAL_INVALID_STOPBITS	 -355
#define HPRT_E_IO_SERIAL_INVALID_FLOWCONTROL -356
/* ethernet port error*/
#define HPRT_E_IO_EHTERNET_CONNECT_ABORT     -361
/* USB port error */
#define HPRT_E_IO_INVALID_USB_PATH	         -371
#define HPRT_E_IO_USB_DEVICE_NOT_FOUND	     -372
#define HPRT_E_IO_USB_DEVICE_BUSY	         -373
/* Extern LIBUSB error */
#define HPRT_E_IO_LIBUSB_E_START	         -1100
#define HPRT_E_IO_LIBUSB_E_END				 -1200
/** Success (no error) */
#define HPRT_E_LIBUSB_SUCCESS                 -1100
/** Input/output error */
#define HPRT_E_LIBUSB_ERROR_IO                -1101
/** Invalid parameter */
#define HPRT_E_LIBUSB_ERROR_INVALID_PARAM     -1102
/** Access denied (insufficient permissions) */
#define HPRT_E_LIBUSB_ERROR_ACCESS            -1103
/** No such device (it may have been disconnected) */
#define HPRT_E_LIBUSB_ERROR_NO_DEVICE         -1104
/** Entity not found */
#define HPRT_E_LIBUSB_ERROR_NOT_FOUND         -1105
/** Resource busy */
#define HPRT_E_LIBUSB_ERROR_BUSY              -1106
/** Operation timed out */
#define HPRT_E_LIBUSB_ERROR_TIMEOUT           -1107
/** Overflow */
#define HPRT_E_LIBUSB_ERROR_OVERFLOW          -1108
/** Pipe error */
#define HPRT_E_LIBUSB_ERROR_PIPE              -1109
/** System call interrupted (perhaps due to signal) */
#define HPRT_E_LIBUSB_ERROR_INTERRUPTED       -1110
/** Insufficient memory */
#define HPRT_E_LIBUSB_ERROR_NO_MEM            -1111
/** Operation not supported or unimplemented on this platform */
#define HPRT_E_LIBUSB_ERROR_NOT_SUPPORTED     -1112

/** Other error */
#define HPRT_E_LIBUSB_ERROR_OTHER   -1199

//====================================

// msr track
#define HPRT_E_MSR_TRACK_NOT_READY			 -401
// smard card
#define HPRT_E_SMART_CARD_NOT_READY			 -411
//encrypt head
#define HPRT_E_EH_SET_ERROR				     -501
#define HPRT_E_EH_DECRYPT_ERROR				 -511

//自定义错误码
#define HPRT_E_CUSTOM_START					 -6000
#define HPRT_E_CUSTOM_Exception				(HPRT_E_CUSTOM_START - 0)			//异常
#define HPRT_E_CUSTOM_Success				(HPRT_E_CUSTOM_START - 1)			//成功
#define HPRT_E_CUSTOM_Start					(HPRT_E_CUSTOM_START - 2)			//开始
#define HPRT_E_CUSTOM_Progress				(HPRT_E_CUSTOM_START - 3)			//进度
#define HPRT_E_CUSTOM_Cancel				(HPRT_E_CUSTOM_START - 4)			//取消
#define HPRT_E_CUSTOM_RestartPrinterFailed	(HPRT_E_CUSTOM_START - 5)			//重启打印机失败
#define HPRT_E_CUSTOM_EnterModeFailed		(HPRT_E_CUSTOM_START - 6)			//进入下载模式失败
#define HPRT_E_CUSTOM_FileIsNull			(HPRT_E_CUSTOM_START - 7)			//文件内容为空
#define HPRT_E_CUSTOM_DataIsMissed			(HPRT_E_CUSTOM_START - 8)			//接收数据缺失
#define HPRT_E_CUSTOM_SlaveResponseFailed	(HPRT_E_CUSTOM_START - 9)			//下位机响应失败
#define HPRT_E_CUSTOM_SpaceNotEnough		(HPRT_E_CUSTOM_START - 10)			//空间不足
#define HPRT_E_CUSTOM_FileDataIsInvalid		(HPRT_E_CUSTOM_START - 11)			//文件数据错误
#define HPRT_E_CUSTOM_NotFoundDevice		(HPRT_E_CUSTOM_START - 12)			//设备未找到

//==============================================================
extern int hprt_last_error;

#ifdef LINUX
#	define EXPORT_ERRORLIB extern 
#elif HPRT_ERROR_EXPORTS
#	ifdef __cplusplus
#		define EXPORT_ERRORLIB extern "C" DLL_EXPORT
#	else
#		define EXPORT_ERRORLIB DLL_EXPORT
#	endif  /* __cplusplus */
#else
#	define EXPORT_ERRORLIB
#endif  /* LINUX */

#ifdef __cplusplus
extern "C"
{
#endif

EXPORT_ERRORLIB int  hprt_get_last_error(void);
EXPORT_ERRORLIB void hprt_set_last_error(int error_no);

#if defined(_UNICODE) || defined (WINCE)
#define hprt_format_error hprt_format_error_wchar
#else
#define hprt_format_error hprt_format_error_ansi
#endif

#ifndef WINCE
	EXPORT_ERRORLIB void hprt_format_error_ansi(int error_no, int langid, unsigned char* buf, int pos, int buf_size);
#endif
EXPORT_ERRORLIB void hprt_format_error_wchar(int error_no, int langid, unsigned char* buf, int pos, int buf_size);

#ifdef __cplusplus
}
#endif

#endif	//_HPRT_ERROR_H_