import datetime
from json import load
import os
import sys
import psycopg2
import redis
from dotenv import load_dotenv


my_directory = os.getcwd()
parent_directory = os.path.dirname(my_directory)

# loading environment variables
env= os.path.join(parent_directory,".env")
if os.path.exists(env):
  load_dotenv(env)

# database configurarion variables from .env
db_host = os.getenv("DB_HOST")
db_port = os.getenv("DB_PORT")
db_user = os.getenv("DB_USER")
db_password = os.getenv("DB_PASSWORD")
db_name = os.getenv("DB_NAME")

#redis config variables
redis_host = os.getenv("REDIS_HOST")
redis_port = int(os.getenv("REDIS_PORT"))
redis_db = int(os.getenv("REDIS_DB"))
print(redis_host,redis_port,redis_db)

# print out psycopg2 errors
def print_psycopg2_exception(err):
    # get details about the exception
    err_type, err_obj, traceback = sys.exc_info()

    # get the line number when exception occured
    line_num = traceback.tb_lineno

    # print the connect() error
    print ("\npsycopg2 ERROR:", err, "on line number:", line_num)
    print ("psycopg2 traceback:", traceback, "-- type:", err_type)

    # psycopg2 extensions.Diagnostics object attribute
    print ("\nextensions.Diagnostics:", err.diag)

    # print the pgcode and pgerror exceptions
    print ("pgerror:", err.pgerror)
    print ("pgcode:", err.pgcode, "\n")


#configurations and processes
try:

    #connect to the database
    r = redis.Redis(host=redis_host, port=redis_port, db=0,decode_responses=True)
    connection = psycopg2.connect(user=db_user,
                                  password=db_password,
                                  host=db_host,
                                  port=db_port,
                                  database=db_name)
    print(connection)
    # get cursor to for query
    cursor = connection.cursor()
    print("connected")

    #temporary_user_files_directory directory contains the temporay files to be appended to the user files
    #temporary_user_files_directory directory contains the temporay files to be appended to the user files
    temporary_user_files_directory = parent_directory + "/uploads"
    print("temp")
    
    # loop to continuesly check redis list(tasks queue),and compare it to the 
    # files in the temporary storage,if it matches ,then read its content, and 
    # add it to the users file

    while True:
      # get the items in the task queue, and get the number of elements in it
      task_list = r.lrange("tasks_queue",0,-1)
      items_in_queue = len(task_list)

      # if there is an item or task in the queue,pop the last on the list(this is a blocking operations)
      # the task item is a tuple that contains the name of the list and the name of the file to be processed
      #  get the name of the the file to be processed ,then iterate over every file in the temporary directory
      #  comparing the filename the file,if any of the file matches ,read it ,and add the contents to the 
      # appropriate user(this can be gotten from the filename>>"username_filename.txt"),and delete temporary file
      if items_in_queue >= 1 :
        print(f"this is the tasks list::{task_list}")
        print(f"there are {items_in_queue} tasks in queue")

        task = r.brpop("tasks_queue",timeout=0)
        file_from_task = task[1]
        print(f"task>>{task}")
        #iterate over every file in the temporary directory
        for filename in os.listdir(temporary_user_files_directory):
          if file_from_task == filename:
            # the files come in the form 'username_filename.txt',split the filename at '_'
            # to a list containing the username and filename
            seperated = filename.split("_")
            print(f"seperated user and file>>{seperated}")
            #username from splitting the filename
            user = seperated[0]

            #if the present index of the for loop is a file,read its contents
            user_temporary_file = os.path.join(temporary_user_files_directory,filename)
            if os.path.isfile(user_temporary_file):
              try:
                user_file = open(user_temporary_file,"r")
                contents = user_file.read()
              except:
                print("error could not open or read file")

            #query the database,for the user whose name is gotten from the filename splitted
            cursor.execute("SELECT name,file_path,last_changed,email FROM users WHERE name=%s",(user,))
  
            user_query = cursor.fetchone()
            print(user_query)

            #create a directory where permanent users files would be kept
            # if it does not exist,but if it does skip
            users_permanent_file_directory = os.path.join(parent_directory,"users_files")
            if not os.path.exists(users_permanent_file_directory):
              os.mkdir(users_permanent_file_directory)

            #get the file path from the query,if its ""(empty string),create a file using the user's name
            #and append to it the content of the temporary file(file from temporary directory with users name)
            username=user_query[0] #username from query
            user_file_path=user_query[1] #user file path from query
            if user_file_path == "":
              user_permanent_file = username + ".txt"
              path = os.path.join(users_permanent_file_directory,user_permanent_file)
              with open(path,"x") as f:
                try:
                  f.write("\n"+contents)
                  print("creating user file")
                  print("file contents appended to user file")
                except:
                  print("Could not open or append to user permanent file")

              update_query = """Update users set file_path=%s,last_changed=%s where name=%s"""
              cursor.execute(update_query,(path,datetime.datetime.now(),username))
              connection.commit()
              present_directory = os.getcwd()
              print(f"present_directory::{present_directory}")

            #if filepath exists,apend contents of temporary file to the file
            else:
              with open(user_file_path,"a") as f:
                f.write("\n"+contents)
                print(user_file_path)
                print("file contents appended to user file")

            #delete temporary file
            os.remove(user_temporary_file)
            print("temporary file deleted")
            print("------------------------------------------")

    
#error while connecting to the database
except psycopg2.OperationalError as db_error:
  print_psycopg2_exception(db_error)
  connection = None

#error connecting to redis
except redis.exceptions.ConnectionError as redis_error:
  print("error connecting to redis>>",redis_error)

#close connection when done         
finally:
  if connection:
    cursor.close()
    connection.close()
    print("PostgreSQL connection is closed")
